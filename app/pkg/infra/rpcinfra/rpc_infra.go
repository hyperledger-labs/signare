// Package rpcinfra provides the infrastructure for implementing a Remote Procedure Call (RPC) server.
package rpcinfra

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra/rpcerrors"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	"github.com/gorilla/mux"
)

// RPCRouter are a set of methods to set up an RPC Router
type RPCRouter interface {
	// HandleRPCRequest handles incoming JSON-RPC requests
	HandleRPCRequest(w http.ResponseWriter, r *http.Request)
	// RegisterRPCHandlerFunc adds a new method and its handler to the router
	RegisterRPCHandlerFunc(method string, handler RPCHandler) error
	// RegisterMiddleware registers a collection of RawMiddlewares
	RegisterMiddleware(middleware ...func(handler http.Handler) http.Handler) error
	// RPCHandler returns the JSON-RPC method handler of a DefaultRPCRouter
	RPCHandler(method string) (RPCHandler, *rpcerrors.RPCError)
	// Methods returns the JSON-RPC registered methods of a DefaultRPCRouter.
	Methods() []string
	// Router returns the mux router
	Router() *mux.Router
}

// Router returns the router of a DefaultRPCRouter
func (rpcRouter *DefaultRPCRouter) Router() *mux.Router {
	return rpcRouter.router
}

// RPCHandler returns the JSON-RPC method handler of a DefaultRPCRouter
func (rpcRouter *DefaultRPCRouter) RPCHandler(method string) (RPCHandler, *rpcerrors.RPCError) {
	handler, ok := rpcRouter.rpcHandlers[method]
	if !ok {
		return nil, rpcerrors.NewMethodNotFound()
	}
	return handler, nil
}

// RegisterRPCHandlerFunc adds a new method and its handler to the router
func (rpcRouter *DefaultRPCRouter) RegisterRPCHandlerFunc(method string, handler RPCHandler) error {
	err := validateMethodName(method)
	if err != nil {
		return err
	}
	rpcRouter.rpcHandlers[method] = handler
	return nil
}

// RegisterMiddleware registers a collection of RawMiddlewares
func (rpcRouter *DefaultRPCRouter) RegisterMiddleware(middlewares ...func(handler http.Handler) http.Handler) error {
	muxMiddleWareArr := make([]mux.MiddlewareFunc, len(middlewares))
	for counter, currentMiddleware := range middlewares {
		muxMiddleWareArr[counter] = mux.MiddlewareFunc(currentMiddleware)
	}

	rpcRouter.router.Use(muxMiddleWareArr...)
	return nil
}

// HandleRPCRequest handles incoming JSON-RPC requests.
func (rpcRouter *DefaultRPCRouter) HandleRPCRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rpcRouter.defaultRPCInfraResponseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument))
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.LogEntry(ctx).Errorf("error closing request body: [%+v]", err)
		}
	}(r.Body)

	var rpcRequest RPCRequest
	if err = json.Unmarshal(body, &rpcRequest); err != nil {
		rpcRouter.defaultRPCInfraResponseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument))
		return
	}

	handler, notFoundRPCError := rpcRouter.RPCHandler(rpcRequest.Method)
	if notFoundRPCError != nil {
		logger.LogEntry(ctx).Errorf("RPC method [%s] not registered", rpcRequest.Method)
		rpcRouter.defaultRPCInfraResponseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusInvalidArgument))
		return
	}

	result, rpcError := handler(ctx, rpcRequest)
	if rpcError != nil {
		httpErr := &httpinfra.HTTPError{}
		httpErr.SetOriginalError(rpcError)
		rpcRouter.defaultRPCInfraResponseHandler.HandleErrorResponse(ctx, w, httpErr)
		return
	}

	// A Notification is a Request object without an "id" member.
	// A Request object that is a Notification signifies the Client's lack of interest in the corresponding Response object,
	// and as such no Response object needs to be returned to the client.
	// src: https://www.jsonrpc.org/specification
	if rpcRequest.ID == nil {
		return
	}

	rpcRouter.defaultRPCInfraResponseHandler.HandleSuccessResponse(ctx, w, httpinfra.ResponseInfo{}, result)
}

// Methods returns the JSON-RPC registered methods of a DefaultRPCRouter.
func (rpcRouter *DefaultRPCRouter) Methods() []string {
	methods := make([]string, 0, len(rpcRouter.rpcHandlers))
	for key := range rpcRouter.rpcHandlers {
		methods = append(methods, key)
	}
	return methods
}

var _ RPCRouter = (*DefaultRPCRouter)(nil)

// DefaultRPCRouterOptions to configure the DefaultRPCRouter
type DefaultRPCRouterOptions struct {
	// DefaultRPCInfraResponseHandler handles the RPC responses of the server
	DefaultRPCInfraResponseHandler *DefaultRPCInfraResponseHandler
	// HTTPMetrics is a set of metrics to count forbidden access attempts
	HTTPMetrics httpinfra.HTTPMetrics
}

// ProvideDefaultRPCRouter creates a new DefaultRPCRouter
func ProvideDefaultRPCRouter(options DefaultRPCRouterOptions) *DefaultRPCRouter {
	return &DefaultRPCRouter{
		rpcHandlers:                    make(map[string]RPCHandler),
		router:                         mux.NewRouter(),
		defaultRPCInfraResponseHandler: options.DefaultRPCInfraResponseHandler,
	}
}

// Method names that begin with "rpc." are reserved for system extensions, and MUST NOT be used for anything else.
// src: https://www.jsonrpc.org/specification#extensions
func validateMethodName(method string) error {
	if strings.HasPrefix(method, "rpc.") {
		return errors.Internal().WithMessage("request method has invalid pattern [rpc.*]")
	}
	return nil
}
