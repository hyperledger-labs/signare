package authentication

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextvalidation"
)

// CreateMiddlewareChain creates a new middleware chain for requests
func (m AuthenticationMiddleware) CreateMiddlewareChain() []func(next http.Handler) http.Handler {
	var middleware []func(next http.Handler) http.Handler

	middleware = append(middleware, m.contextDefinition.DefineAction)
	middleware = append(middleware, m.requestContextValidation.ValidateAction)

	middleware = append(middleware, m.contextDefinition.DefineUser)
	middleware = append(middleware, m.requestContextValidation.ValidateUser)

	middleware = append(middleware, m.contextDefinition.DefineApplication)
	middleware = append(middleware, m.requestContextValidation.ValidateApplication)

	return middleware
}

// AuthenticationMiddlewareOptions are the set of fields to create an AuthenticationMiddleware
type AuthenticationMiddlewareOptions struct {
	// ContextDefinition is a set of middleware functions to create middleware chains
	ContextDefinition contextdefinition.ContextDefinition
	// RequestContextValidation defines a middleware that checks if the context contains the expected auth values
	RequestContextValidation *contextvalidation.RequestContextValidation
}

// AuthenticationMiddleware is the middleware used for authenticate the user making a request
type AuthenticationMiddleware struct {
	contextDefinition        contextdefinition.ContextDefinition
	requestContextValidation *contextvalidation.RequestContextValidation
}

// ProvideAuthenticationMiddleware provides an instance of an AuthenticationMiddleware
func ProvideAuthenticationMiddleware(options AuthenticationMiddlewareOptions) (*AuthenticationMiddleware, error) {
	if options.ContextDefinition == nil {
		return nil, errors.New("mandatory 'ContextDefinition' not provided")
	}
	if options.RequestContextValidation == nil {
		return nil, errors.New("mandatory 'RequestContextValidation' not provided")
	}

	return &AuthenticationMiddleware{
		contextDefinition:        options.ContextDefinition,
		requestContextValidation: options.RequestContextValidation,
	}, nil
}
