// Package httpinfra provides infrastructure for managing HTTP requests offering a structured approach
// to handling incoming HTTP requests and routing them to the appropriate handlers
package httpinfra

import (
	"net/http"

	"github.com/gorilla/mux"
)

// HTTPRouter are a set of methods to set up an HTTP Router
type HTTPRouter interface {
	// RegisterRawHandlerWithMatcherFunc registers a route using HandlerMatcherFuncOptions and RawHandler
	RegisterRawHandlerWithMatcherFunc(matchOptions HandlerMatcherFuncOptions, handler RawHandler) error
	// RegisterRawHandler registers a route using HandlerMatchOptions and RawHandler
	RegisterRawHandler(matchOptions HandlerMatchOptions, handler RawHandler) error
	// RegisterHandlerFunc registers a route using HandlerMatchOptions and RawHandler
	RegisterHandlerFunc(matchOptions HandlerMatchOptions, handler http.Handler) error
	// RegisterMiddleware registers a collection of RawMiddlewares
	RegisterMiddleware(middleware ...func(handler http.Handler) http.Handler) error
	// Router returns the mux router
	Router() *mux.Router
}

// MainRouter returns the router of a DefaultHTTPRouter
func (httpRouter *DefaultHTTPRouter) MainRouter() *mux.Router {
	return httpRouter.router
}

// RegisterRawHandlerWithMatcherFunc registers a route using HandlerMatcherFuncOptions and RawHandler
func (httpRouter *DefaultHTTPRouter) RegisterRawHandlerWithMatcherFunc(matchOptions HandlerMatcherFuncOptions, handler RawHandler) error {
	httpRouter.router.MatcherFunc(matchOptions.MatcherFunc).HandlerFunc(handler).Methods(matchOptions.Methods...)
	return nil
}

// RegisterRawHandler registers a route using HandlerMatchOptions and RawHandler
func (httpRouter *DefaultHTTPRouter) RegisterRawHandler(matchOptions HandlerMatchOptions, handler RawHandler) error {
	httpRouter.router.HandleFunc(matchOptions.Path, handler).Methods(matchOptions.Methods...).Name(matchOptions.Action)
	return nil
}

// RegisterHandlerFunc registers a route using HandlerMatchOptions and RawHandler
func (httpRouter *DefaultHTTPRouter) RegisterHandlerFunc(matchOptions HandlerMatchOptions, handler http.Handler) error {
	httpRouter.router.Handle(matchOptions.Path, handler).Methods(matchOptions.Methods...).Name(matchOptions.Action)
	return nil
}

// RegisterMiddleware registers a collection of RawMiddlewares
func (httpRouter *DefaultHTTPRouter) RegisterMiddleware(middleware ...func(http.Handler) http.Handler) error {
	muxMiddleWareArr := make([]mux.MiddlewareFunc, len(middleware))
	for counter, currentMiddleware := range middleware {
		muxMiddleWareArr[counter] = currentMiddleware
	}

	httpRouter.router.Use(muxMiddleWareArr...)
	return nil
}

// Router returns the mux router
func (httpRouter *DefaultHTTPRouter) Router() *mux.Router {
	return httpRouter.router
}

// DefaultHTTPRouter holds the router to handle HTTP incoming connections.
type DefaultHTTPRouter struct {
	router *mux.Router
}

// ProvideHTTPRouter creates a DefaultHTTPRouter.
func ProvideHTTPRouter() *DefaultHTTPRouter {
	router := mux.NewRouter()
	return &DefaultHTTPRouter{
		router: router,
	}
}
