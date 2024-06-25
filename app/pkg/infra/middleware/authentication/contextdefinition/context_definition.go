package contextdefinition

import (
	"net/http"
)

// ContextDefinition is a set of middleware functions to create middleware chains
type ContextDefinition interface {
	// DefineUser defines the user within the context of the request
	DefineUser(next http.Handler) http.Handler
	// DefineApplication defines the application within the context of the request
	DefineApplication(next http.Handler) http.Handler
	// DefineAction defines the action within the context of the request
	DefineAction(next http.Handler) http.Handler
}
