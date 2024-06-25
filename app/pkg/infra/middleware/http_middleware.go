package middleware

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/telemetry"
)

// Create the middlewares chains for HTTP
func (f HTTPMiddlewareFactory) Create() []func(handler http.Handler) http.Handler {
	fullChain := make([]func(handler http.Handler) http.Handler, 0)

	telemetryMiddleware := f.telemetryMiddleware.CreateMiddlewareChain()
	authenticationMiddlewareChain := f.authenticationMiddleware.CreateMiddlewareChain()
	authorizationMiddleware := f.authorizationMiddleware.CreateMiddlewareChain(false)

	fullChain = append(fullChain, telemetryMiddleware...)
	fullChain = append(fullChain, authenticationMiddlewareChain...)
	fullChain = append(fullChain, authorizationMiddleware...)

	return fullChain
}

// HTTPMiddlewareFactory creates middlewares chains for HTTP
type HTTPMiddlewareFactory struct {
	authenticationMiddleware *authentication.AuthenticationMiddleware
	authorizationMiddleware  *authorization.AuthorizationMiddleware
	telemetryMiddleware      *telemetry.TelemetryMiddleware
}

// HTTPMiddlewareFactoryOptions are the set of fields to create an HTTPMiddlewareFactory
type HTTPMiddlewareFactoryOptions struct {
	// AuthenticationMiddleware is the middleware used for authenticate the user making a request
	AuthenticationMiddleware *authentication.AuthenticationMiddleware
	// AuthorizationMiddleware is the middleware used for authorize the user making a request
	AuthorizationMiddleware *authorization.AuthorizationMiddleware
	// TelemetryMiddleware is the middleware used for handling telemetry within requests
	TelemetryMiddleware *telemetry.TelemetryMiddleware
}

// ProvideHTTPMiddlewareFactory provides an instance of an HTTPMiddlewareFactory
func ProvideHTTPMiddlewareFactory(options HTTPMiddlewareFactoryOptions) (*HTTPMiddlewareFactory, error) {
	if options.AuthorizationMiddleware == nil {
		return nil, errors.New("mandatory 'AuthorizationMiddleware' not provided")
	}
	if options.AuthenticationMiddleware == nil {
		return nil, errors.New("mandatory 'AuthenticationMiddleware' not provided")
	}
	if options.TelemetryMiddleware == nil {
		return nil, errors.New("mandatory 'TelemetryMiddleware' not provided")
	}
	return &HTTPMiddlewareFactory{
		authenticationMiddleware: options.AuthenticationMiddleware,
		authorizationMiddleware:  options.AuthorizationMiddleware,
		telemetryMiddleware:      options.TelemetryMiddleware,
	}, nil
}
