package telemetry

import (
	"errors"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/infra/middleware/telemetry/tracer"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace/noop"
)

// CreateMiddlewareChain creates a new middleware chain for requests
func (m TelemetryMiddleware) CreateMiddlewareChain() []func(next http.Handler) http.Handler {
	var middleware []func(next http.Handler) http.Handler

	tracerPropagator := propagation.TraceContext{}
	tracerProvider := noop.NewTracerProvider()

	// Send empty service so the request host is used to determine the server name.
	middleware = append(middleware, otelmux.Middleware("",
		otelmux.WithPropagators(tracerPropagator),
		otelmux.WithTracerProvider(tracerProvider),
	))

	middleware = append(middleware, m.httpContextTracer.HandleTracing)

	return middleware
}

// TelemetryMiddlewareOptions are the set of fields to create an TelemetryMiddleware
type TelemetryMiddlewareOptions struct {
	HTTPContextTracer *tracer.HTTPContextTracer
}

// TelemetryMiddleware is the middleware used for managing telemetry within requests.
type TelemetryMiddleware struct {
	httpContextTracer *tracer.HTTPContextTracer
}

// ProvideTelemetryMiddleware provides an instance of an TelemetryMiddleware
func ProvideTelemetryMiddleware(options TelemetryMiddlewareOptions) (*TelemetryMiddleware, error) {
	if options.HTTPContextTracer == nil {
		return nil, errors.New("mandatory 'HTTPContextTracer' not provided")
	}
	return &TelemetryMiddleware{
		httpContextTracer: options.HTTPContextTracer,
	}, nil
}
