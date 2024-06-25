package authorization

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authorization/pep"
)

// CreateMiddlewareChain creates a new middleware chain for requests
func (m AuthorizationMiddleware) CreateMiddlewareChain(enableAccountUserAuthorization bool) []func(next http.Handler) http.Handler {
	var middleware []func(next http.Handler) http.Handler

	middleware = append(middleware, m.httpPolicyEnforcementPoint.AuthorizeUser)

	if enableAccountUserAuthorization {
		middleware = append(middleware, m.rpcPolicyEnforcementPoint.AuthorizeAccount)
	}

	return middleware
}

// AuthorizationMiddlewareOptions are the set of fields to create an AuthorizationMiddleware
type AuthorizationMiddlewareOptions struct {
	// HTTPPolicyEnforcementPoint checks if the user is authorized to perform a given action
	HTTPPolicyEnforcementPoint *pep.HTTPPolicyEnforcementPoint
	// RPCPolicyEnforcementPoint checks if the user is authorized to perform a given action
	RPCPolicyEnforcementPoint *pep.RPCPolicyEnforcementPoint
}

// AuthorizationMiddleware is the middleware used for authorize the user making a request
type AuthorizationMiddleware struct {
	httpPolicyEnforcementPoint *pep.HTTPPolicyEnforcementPoint
	rpcPolicyEnforcementPoint  *pep.RPCPolicyEnforcementPoint
}

// ProvideAuthorizationMiddleware provides an instance of an AuthorizationMiddleware
func ProvideAuthorizationMiddleware(options AuthorizationMiddlewareOptions) (*AuthorizationMiddleware, error) {
	if options.HTTPPolicyEnforcementPoint == nil {
		return nil, errors.New("'httpPolicyEnforcementPoint' field is mandatory and cannot be empty")
	}
	if options.RPCPolicyEnforcementPoint == nil {
		return nil, errors.New("'RPCPolicyEnforcementPoint' field is mandatory and cannot be empty")
	}
	return &AuthorizationMiddleware{
		httpPolicyEnforcementPoint: options.HTTPPolicyEnforcementPoint,
		rpcPolicyEnforcementPoint:  options.RPCPolicyEnforcementPoint,
	}, nil
}
