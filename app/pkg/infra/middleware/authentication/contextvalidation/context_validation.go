package contextvalidation

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"

	"github.com/gorilla/mux"
)

// ValidateUser returns middleware that checks if the context contains a non-empty user
func (m *RequestContextValidation) ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		validated := false

		user, err := requestcontext.UserFromContext(ctx)
		if err != nil {
			m.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		if user != nil && len(*user) > 0 {
			validated = true
		}

		if !validated {
			m.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusPermissionDenied))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateApplication returns middleware that checks if the context contains a non-empty application
func (m *RequestContextValidation) ValidateApplication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// the error from the method below is ignored because the middleware allows requests without application header
		application, err := requestcontext.ApplicationFromContext(ctx)
		if err != nil {
			// if there is not an application in the context, then this means that the user could be a signer administrator
			next.ServeHTTP(w, r)
			return
		}

		params := mux.Vars(r)
		applicationPathParam := params["applicationId"]

		// JSON-RPC endpoints and /admin endpoints don't need access control verification
		if len(applicationPathParam) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		applicationValidated := application != nil && *application == applicationPathParam
		if !applicationValidated {
			m.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusPermissionDenied))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateAction returns middleware that checks if the context contains a non-empty action
func (m *RequestContextValidation) ValidateAction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		validated := false

		action, err := requestcontext.ActionFromContext(ctx)
		if err != nil {
			m.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPErrorFromError(ctx, err, httpinfra.StatusPermissionDenied))
			return
		}
		if action != nil && len(*action) > 0 {
			validated = true
		}

		if !validated {
			m.responseHandler.HandleErrorResponse(ctx, w, httpinfra.NewHTTPError(httpinfra.StatusPermissionDenied))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequestContextValidation defines a middleware that checks if the context contains the expected auth values
type RequestContextValidation struct {
	responseHandler httpinfra.HTTPResponseHandler
}

// RequestContextValidationOptions options to create a new RequestContextValidation
type RequestContextValidationOptions struct {
	// ResponseHandler exposes functionality to handle HTTP responses
	ResponseHandler httpinfra.HTTPResponseHandler
}

// ProvideRequestContextValidation returns a RequestContextValidation with the given options
func ProvideRequestContextValidation(options RequestContextValidationOptions) (*RequestContextValidation, error) {
	if options.ResponseHandler == nil {
		return nil, errors.New("mandatory 'DefaultRPCInfraResponseHandler' not provided")
	}
	return &RequestContextValidation{
		responseHandler: options.ResponseHandler,
	}, nil
}
