package middleware

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/entrypoint/rpcbatchrequestsupport"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/telemetry"
)

// Create the middlewares chains for the RPC protocol
func (f RPCMiddlewareFactory) Create() []func(handler http.Handler) http.Handler {
	fullChain := make([]func(handler http.Handler) http.Handler, 0)

	telemetryMiddleware := f.telemetryMiddleware.CreateMiddlewareChain()
	rpcBatchRequestSupportMiddlewareChain := f.rpcBatchRequestSupport.FanOutRPCBatchRequest
	authenticationMiddlewareChain := f.authenticationMiddleware.CreateMiddlewareChain()
	authorizationMiddleware := f.authorizationMiddleware.CreateMiddlewareChain(true)

	fullChain = append(fullChain, telemetryMiddleware...)
	fullChain = append(fullChain, rpcBatchRequestSupportMiddlewareChain)
	fullChain = append(fullChain, authenticationMiddlewareChain...)
	fullChain = append(fullChain, authorizationMiddleware...)

	return fullChain
}

// RPCMiddlewareFactory creates middlewares chains for the RPC protocol
type RPCMiddlewareFactory struct {
	authenticationMiddleware *authentication.AuthenticationMiddleware
	authorizationMiddleware  *authorization.AuthorizationMiddleware
	rpcBatchRequestSupport   *rpcbatchrequestsupport.RPCBatchRequestSupportMiddleware
	telemetryMiddleware      *telemetry.TelemetryMiddleware
}

// RPCMiddlewareFactoryOptions are the set of fields to create an RPCMiddlewareFactory
type RPCMiddlewareFactoryOptions struct {
	// AuthenticationMiddleware is the middleware used for authenticate the user making a request
	AuthenticationMiddleware *authentication.AuthenticationMiddleware
	// AuthorizationMiddleware is the middleware used for authorize the user making a request
	AuthorizationMiddleware *authorization.AuthorizationMiddleware
	// RPCBatchRequestSupportMiddleware enables support for batch requests as defined in https://www.jsonrpc.org/specification#batch
	RPCBatchRequestSupportMiddleware *rpcbatchrequestsupport.RPCBatchRequestSupportMiddleware
	// TelemetryMiddleware is the middleware used for handling telemetry within requests
	TelemetryMiddleware *telemetry.TelemetryMiddleware
}

// ProvideRPCMiddlewareFactory provides an instance of an RPCMiddlewareFactory
func ProvideRPCMiddlewareFactory(options RPCMiddlewareFactoryOptions) (*RPCMiddlewareFactory, error) {
	if options.RPCBatchRequestSupportMiddleware == nil {
		return nil, errors.New("mandatory 'RPCBatchRequestSupportMiddleware' not provided")
	}
	if options.AuthorizationMiddleware == nil {
		return nil, errors.New("mandatory 'AuthorizationMiddleware' not provided")
	}
	if options.AuthenticationMiddleware == nil {
		return nil, errors.New("mandatory 'AuthenticationMiddleware' not provided")
	}
	if options.TelemetryMiddleware == nil {
		return nil, errors.New("mandatory 'TelemetryMiddleware' not provided")
	}
	return &RPCMiddlewareFactory{
		authenticationMiddleware: options.AuthenticationMiddleware,
		authorizationMiddleware:  options.AuthorizationMiddleware,
		rpcBatchRequestSupport:   options.RPCBatchRequestSupportMiddleware,
		telemetryMiddleware:      options.TelemetryMiddleware,
	}, nil
}
