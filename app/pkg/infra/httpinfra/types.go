package httpinfra

import (
	"net/http"

	"github.com/gorilla/mux"
)

// ErrorResponseWrapper wraps error for ErrorResponse
type ErrorResponseWrapper struct {
	// Error wrapped error
	Error ErrorResponse `json:"error"`
}

// ErrorResponse defines a response for a failed http call
type ErrorResponse struct {
	// Code error code
	Code int `json:"code"`
	// Status error status
	Status string `json:"status"`
	// Details error details
	Details ErrorResponseDetails `json:"details"`
}

// ErrorResponseDetails details a error
type ErrorResponseDetails struct {
	// TraceableErrorId traceable error identifier
	TraceableErrorId string `json:"traceableErrorId"`
	// Message string message
	Message string `json:"message"`
}

// HandlerMatcherFuncOptions provides information to register a new route with a custom matcher for a path from an HTTP request with specific HTTP methods
type HandlerMatcherFuncOptions struct {
	// MatcherFunc function that checks if an HTTP request path is handled by the route
	MatcherFunc func(*http.Request, *mux.RouteMatch) bool
	// Methods http methods for the route
	Methods []string
}

// HandlerMatchOptions configures a handler to register a new route with specific HTTP methods
type HandlerMatchOptions struct {
	// Path that is handled by the route
	Path string
	// Methods http methods for the route
	Methods []string
	// Action configures the action to be assigned as the name of the route
	Action string
}

// RawHandler function that handles a http request to emit an HTTP response
type RawHandler func(w http.ResponseWriter, r *http.Request)
