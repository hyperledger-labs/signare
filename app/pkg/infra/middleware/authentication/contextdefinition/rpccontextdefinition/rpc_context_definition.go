package rpccontextdefinition

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
	"github.com/hyperledger-labs/signare/app/pkg/infra/rpcinfra"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/utils"
)

const (
	supportedRPCVersion = "2.0"
)

// DefineUser defines the user within the context of the request
func (middleware *RPCContextDefinition) DefineUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := contextdefinition.DefaultUserHeader
		if middleware.authHeadersConfiguration.UserRequestHeader != "" {
			headerKey = middleware.authHeadersConfiguration.UserRequestHeader
		}
		userID := r.Header.Get(headerKey)

		trimmedUser := strings.TrimSpace(userID)
		ctx := context.WithValue(r.Context(), requestcontext.UserContextKey, trimmedUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DefineApplication defines the application within the context of the request
func (middleware *RPCContextDefinition) DefineApplication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := contextdefinition.DefaultApplicationHeader
		if middleware.authHeadersConfiguration.ApplicationRequestHeader != "" {
			headerKey = middleware.authHeadersConfiguration.ApplicationRequestHeader
		}
		applicationID := r.Header.Get(headerKey)

		trimmedApplication := strings.TrimSpace(applicationID)
		ctx := context.WithValue(r.Context(), requestcontext.ApplicationContextKey, trimmedApplication)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DefineAction defines the action within the context of the request
func (middleware *RPCContextDefinition) DefineAction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var match mux.RouteMatch
		matches := middleware.router.Router().Match(r, &match)
		if !matches || match.MatchErr != nil {
			return
		}

		actionID := match.Route.GetName()

		var rpcRequest RPCRequest
		err := utils.ReadAndResetCloser(&r.Body, &rpcRequest)
		if err != nil {
			middleware.responseHandler.HandleErrorResponse(r.Context(), w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}

		if rpcRequest.RPCVersion != supportedRPCVersion {
			logger.LogEntry(r.Context()).Errorf("request parameter [jsonrpc] must be exactly '%s'", supportedRPCVersion)
			middleware.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusPermissionDenied))
			return
		}

		composedActionID := fmt.Sprintf("%s.%s", actionID, rpcRequest.Method)
		ctx = context.WithValue(ctx, requestcontext.ActionContextKey, composedActionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RPCContextDefinitionOptions configures RPCContextDefinition
type RPCContextDefinitionOptions struct {
	// AuthHeadersConfiguration configure the auth headers of the requests
	AuthHeadersConfiguration contextdefinition.AuthHeadersConfiguration
	// ResponseHandler exposes functionality to handle HTTP responses
	ResponseHandler httpinfra.HTTPResponseHandler
	// RPCRouter are a set of methods to set up an RPC Router
	RPCRouter rpcinfra.RPCRouter
}

// RPCContextDefinition defines authorization configuration for an application and a user
type RPCContextDefinition struct {
	authHeadersConfiguration contextdefinition.AuthHeadersConfiguration
	responseHandler          httpinfra.HTTPResponseHandler
	router                   rpcinfra.RPCRouter
}

// ProvideRPCContextDefinitionFromHeaders returns RPCContextDefinition with the given options
func ProvideRPCContextDefinitionFromHeaders(options RPCContextDefinitionOptions) (*RPCContextDefinition, error) {
	if options.RPCRouter == nil {
		return nil, errors.Internal().WithMessage("'HTTPRouter' field is mandatory")
	}
	if options.ResponseHandler == nil {
		return nil, errors.Internal().WithMessage("'DefaultRPCInfraResponseHandler' field is mandatory")
	}
	return &RPCContextDefinition{
		authHeadersConfiguration: options.AuthHeadersConfiguration,
		responseHandler:          options.ResponseHandler,
		router:                   options.RPCRouter,
	}, nil
}
