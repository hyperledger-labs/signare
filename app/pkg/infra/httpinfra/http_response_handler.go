package httpinfra

import (
	"context"
	"encoding/json"
	"net/http"
)

// HTTPResponseHandler exposes functionality to handle HTTP responses
type HTTPResponseHandler interface {
	// HandleErrorResponse handles the HTTP response when the operation returned an error
	HandleErrorResponse(ctx context.Context, w http.ResponseWriter, receivedError *HTTPError)
	// HandleSuccessResponse handles the HTTP response when the operation succeeded
	HandleSuccessResponse(ctx context.Context, w http.ResponseWriter, responseInfo ResponseInfo, responseData interface{})
}

// ValidationResult result of a validation
type ValidationResult struct {
	// Valid is true if the result is valid
	Valid bool
	// ValidReason is the reason why Valid is false
	NotValidReason string
}

// ResponseType represents the response type
type ResponseType string

// ResponseInfo response information
type ResponseInfo struct {
	// ResponseType represents the response type
	ResponseType
}

const (
	// ResponseTypeOk successful response
	ResponseTypeOk ResponseType = "Ok"
	// ResponseTypeCreated created response
	ResponseTypeCreated ResponseType = "Created"
	// ResponseTypeAccepted accepted response
	ResponseTypeAccepted ResponseType = "Accepted"
)

// Optional representation of a container object that may or may not contain a value
type Optional[T any] struct {
	isDefined bool
	value     *T
}

// UnmarshalJSON unmarshal JSON data into the pointer struct
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	o.isDefined = true
	return json.Unmarshal(data, &o.value)
}

// IsEmpty returns true if the optional does not contain a value
func (o *Optional[T]) IsEmpty() bool {
	return !o.isDefined
}

// GetValue returns the value
func (o *Optional[T]) GetValue() *T {
	return o.value
}

// SetValue sets the contained value
func (o *Optional[T]) SetValue(updatedValue *T) {
	o.isDefined = true
	o.value = updatedValue
}
