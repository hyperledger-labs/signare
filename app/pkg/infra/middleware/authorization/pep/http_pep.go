package pep

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
)

// AuthorizeUser filters if the user making the request has authorization to perform the action
// nolint: staticcheck
func (policyEnforcementPoint *HTTPPolicyEnforcementPoint) AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, err := requestcontext.UserFromContext(ctx)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		actionID, err := requestcontext.ActionFromContext(ctx)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		// We ignore the error since signer-admins won't use the application header.
		application, _ := requestcontext.ApplicationFromContext(ctx)

		authorizeUserInput := AuthorizeUserInput{
			UserID:        *user,
			ApplicationID: application,
			ActionID:      *actionID,
		}

		_, err = policyEnforcementPoint.userPolicyDecisionPointAdapter.AuthorizeUser(ctx, authorizeUserInput)
		if err != nil {
			policyEnforcementPoint.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// HTTPPolicyEnforcementPointOptions are the set of fields to create an HTTPPolicyEnforcementPoint
type HTTPPolicyEnforcementPointOptions struct {
	// ResponseHandler exposes functionality to handle HTTP responses
	ResponseHandler httpinfra.HTTPResponseHandler
	// UserPolicyDecisionPointPort is a port to adapt authorization checks
	UserPolicyDecisionPointAdapter UserPolicyDecisionPointPort
	// AccountUserPolicyDecisionPointPort is a port to adapt account usage authorization checks
	AccountUserPolicyDecisionPointAdapter AccountUserPolicyDecisionPointPort
}

// HTTPPolicyEnforcementPoint is a set of methods to check user authorization
type HTTPPolicyEnforcementPoint struct {
	responseHandler                       httpinfra.HTTPResponseHandler
	userPolicyDecisionPointAdapter        UserPolicyDecisionPointPort
	accountUserPolicyDecisionPointAdapter AccountUserPolicyDecisionPointPort
}

// ProvideHTTPPolicyEnforcementPoint provides an instance of an HTTPPolicyEnforcementPoint
func ProvideHTTPPolicyEnforcementPoint(options HTTPPolicyEnforcementPointOptions) (*HTTPPolicyEnforcementPoint, error) {
	if options.ResponseHandler == nil {
		return nil, errors.New("mandatory 'ResponseHandler' not provided")
	}
	if options.UserPolicyDecisionPointAdapter == nil {
		return nil, errors.New("mandatory 'UserPolicyDecisionPointAdapter' not provided")
	}
	if options.AccountUserPolicyDecisionPointAdapter == nil {
		return nil, errors.New("mandatory 'AccountUserPolicyDecisionPointAdapter' not provided")
	}

	return &HTTPPolicyEnforcementPoint{
		responseHandler:                       options.ResponseHandler,
		userPolicyDecisionPointAdapter:        options.UserPolicyDecisionPointAdapter,
		accountUserPolicyDecisionPointAdapter: options.AccountUserPolicyDecisionPointAdapter,
	}, nil
}
