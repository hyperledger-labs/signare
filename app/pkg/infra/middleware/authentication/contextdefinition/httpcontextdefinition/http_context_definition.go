package httpcontextdefinition

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/authentication/contextdefinition"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
)

// DefineUser defines the user within the context of the request
func (m *HTTPContextDefinition) DefineUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := contextdefinition.DefaultUserHeader
		if m.authHeadersConfiguration.UserRequestHeader != "" {
			headerKey = m.authHeadersConfiguration.UserRequestHeader
		}
		userID := r.Header.Get(headerKey)

		trimmedUser := strings.TrimSpace(userID)
		ctx := context.WithValue(r.Context(), requestcontext.UserContextKey, trimmedUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DefineApplication defines the application within the context of the request
func (m *HTTPContextDefinition) DefineApplication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerKey := contextdefinition.DefaultApplicationHeader
		if m.authHeadersConfiguration.ApplicationRequestHeader != "" {
			headerKey = m.authHeadersConfiguration.ApplicationRequestHeader
		}
		applicationID := r.Header.Get(headerKey)
		if applicationID == "" {
			next.ServeHTTP(w, r)
			return
		}

		trimmedApplicationID := strings.TrimSpace(applicationID)
		ctx := context.WithValue(r.Context(), requestcontext.ApplicationContextKey, trimmedApplicationID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// DefineAction defines the action within the context of the request
func (m *HTTPContextDefinition) DefineAction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var match mux.RouteMatch
		matches := m.httpRouter.Router().Match(r, &match)
		if !matches || match.MatchErr != nil {
			return
		}

		actionID := match.Route.GetName()
		ctx := context.WithValue(r.Context(), requestcontext.ActionContextKey, actionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HTTPContextDefinitionOptions configures HTTPContextDefinition
type HTTPContextDefinitionOptions struct {
	// AuthHeadersConfiguration configure the auth headers of the requests
	AuthHeadersConfiguration contextdefinition.AuthHeadersConfiguration
	// HTTPRouter are a set of methods to set up an HTTP HTTPRouter
	HTTPRouter httpinfra.HTTPRouter
}

// HTTPContextDefinition defines authorization configuration for an application and a user
type HTTPContextDefinition struct {
	authHeadersConfiguration contextdefinition.AuthHeadersConfiguration
	httpRouter               httpinfra.HTTPRouter
}

// ProvideHTTPContextDefinition returns HTTPContextDefinition with the given options
func ProvideHTTPContextDefinition(options HTTPContextDefinitionOptions) (*HTTPContextDefinition, error) {
	if options.HTTPRouter == nil {
		return nil, errors.New("'HTTPRouter' field is mandatory")
	}
	return &HTTPContextDefinition{
		authHeadersConfiguration: options.AuthHeadersConfiguration,
		httpRouter:               options.HTTPRouter,
	}, nil
}
