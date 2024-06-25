// Package requestcontext defines a utility to enhance context information from HTTP request information
package requestcontext

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

const (
	UserContextKey               entities.ContextKey = "X-Auth-User-Id"
	ApplicationContextKey        entities.ContextKey = "X-Auth-Application-Id"
	ActionContextKey             entities.ContextKey = "X-Action-Id"
	RPCRequestIDKey              entities.ContextKey = "X-RPC-Request-Id"
	TraceParentTraceIDContextKey entities.ContextKey = "traceparent.trace.id"
	TraceParentSpanIDContextKey  entities.ContextKey = "traceparent.span.id"
)

// retrieveStringValueFromContextKey returns the value as string of a context key and returns an error if the key is not present in the context
func retrieveStringValueFromContextKey(ctx context.Context, key entities.ContextKey) (*string, error) {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return nil, fmt.Errorf("%s not found in context", key)
	}

	return &value, nil
}

// retrieveValueFromContextKey returns the value of a context key and returns an error if the key is not present in the context
func retrieveValueFromContextKey(ctx context.Context, key entities.ContextKey) (*any, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, fmt.Errorf("%s not found in context", key)
	}
	return &value, nil
}

// UserFromContext returns user from context or an error if it fails
func UserFromContext(ctx context.Context) (*string, error) {
	return retrieveStringValueFromContextKey(ctx, UserContextKey)
}

// ApplicationFromContext returns application from context or an error if it fails
func ApplicationFromContext(ctx context.Context) (*string, error) {
	return retrieveStringValueFromContextKey(ctx, ApplicationContextKey)
}

// ActionFromContext returns actionID from context or an error if it fails
func ActionFromContext(ctx context.Context) (*string, error) {
	return retrieveStringValueFromContextKey(ctx, ActionContextKey)
}

// RPCRequestIDFromContext returns requestID from context or an error if it fails
func RPCRequestIDFromContext(ctx context.Context) (*any, error) {
	return retrieveValueFromContextKey(ctx, RPCRequestIDKey)
}

// SpanIDFromContext returns spanID from context or an error if it fails
func SpanIDFromContext(ctx context.Context) (*string, error) {
	return retrieveStringValueFromContextKey(ctx, TraceParentSpanIDContextKey)
}

// AuthConfig information on the user and application authorized in the http request as headers
type AuthConfig struct {
	// UserRequestHeader is the key in the header for the user id
	UserRequestHeader entities.ContextKey
	// ApplicationRequestHeader is the key in the header for the application id
	ApplicationRequestHeader entities.ContextKey
}
