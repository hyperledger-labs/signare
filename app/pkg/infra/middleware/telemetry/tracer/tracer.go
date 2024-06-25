package tracer

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"

	"go.opentelemetry.io/otel/trace"
)

// HandleTracing parses, validates and propagates tracing headers as specified in the trace context specification.
// More information about the standard: https://w3.org/TR/trace-context
func (m *HTTPContextTracer) HandleTracing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rawTraceParent := r.Header.Get("traceparent")
		traceParent := trace.SpanContextFromContext(ctx)

		version, validVersion := parseVersion(rawTraceParent)
		if !validVersion {
			newTraceParent, err := resetTraceContext(traceParent)
			if err != nil {
				logger.LogEntry(ctx).Errorf("failed invalidating span context: %v", err)
				next.ServeHTTP(w, r)
				return
			}
			traceParent = *newTraceParent
		}

		if !traceParent.IsValid() {
			newTraceParent, err := resetTraceContext(traceParent)
			if err != nil {
				logger.LogEntry(ctx).Errorf("failed invalidating span context: %v", err)
				next.ServeHTTP(w, r)
				return
			}
			traceParent = *newTraceParent
		}

		contextTraceState := traceParent.TraceState().String()
		contextTraceParent := fmt.Sprintf("%s-%s-%s-%s",
			version.String(),
			traceParent.TraceID().String(),
			traceParent.SpanID().String(),
			traceParent.TraceFlags().String(),
		)

		// add trace context headers to response only when responding to systems that participated in trace
		if traceParent.IsRemote() {
			w.Header().Add(traceParentHeader, contextTraceParent)
			w.Header().Add(traceStateHeader, contextTraceState)
		}

		ctx = context.WithValue(ctx, requestcontext.TraceParentTraceIDContextKey, traceParent.TraceID().String())
		ctx = context.WithValue(ctx, requestcontext.TraceParentSpanIDContextKey, traceParent.SpanID().String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseVersion attempts to parse the trace parent version. If it succeeds, it returns the parsed version and the flag is set to true.
// If it fails, it returns the current max version supported and the flag is set to false.
func parseVersion(rawTraceParent string) (contextVersion, bool) {
	if len(rawTraceParent) != traceParentLength {
		return traceContextVersion, false
	}

	version, err := strconv.Atoi(rawTraceParent[:2])
	if err != nil {
		return traceContextVersion, false
	}
	// version must be followed by a hyphen
	if rawTraceParent[2:3] != "-" {
		return traceContextVersion, false
	}
	v := contextVersion(version)

	return v, true
}

func resetTraceContext(span trace.SpanContext) (*trace.SpanContext, error) {
	traceID, err := newTraceID()
	if err != nil {
		return nil, err
	}
	spanID, err := newSpanID()
	if err != nil {
		return nil, err
	}
	span = span.WithTraceID(*traceID)
	span = span.WithSpanID(*spanID)
	span = span.WithTraceState(trace.TraceState{})
	span = span.WithRemote(false)
	return &span, nil
}

func newTraceID() (*trace.TraceID, error) {
	randID, err := randomHexString(traceIDLength)
	if err != nil {
		return nil, err
	}
	traceID, err := trace.TraceIDFromHex(randID)
	if err != nil {
		return nil, err
	}
	return &traceID, nil
}

func newSpanID() (*trace.SpanID, error) {
	randID, err := randomHexString(spanIDLength)
	if err != nil {
		return nil, err
	}
	spanID, err := trace.SpanIDFromHex(randID)
	if err != nil {
		return nil, err
	}
	return &spanID, nil
}

func randomHexString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// HTTPContextTracer defines a middleware that handles tracing within requests.
type HTTPContextTracer struct {
}

// HTTPContextTracerOptions options to create a new HTTPContextTracer
type HTTPContextTracerOptions struct {
}

// ProvideHTTPContextTracer returns a HTTPContextTracer with the given options
func ProvideHTTPContextTracer(_ HTTPContextTracerOptions) (*HTTPContextTracer, error) {
	return &HTTPContextTracer{}, nil
}
