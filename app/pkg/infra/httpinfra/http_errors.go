package httpinfra

import (
	"context"
	"net/http"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	"github.com/google/uuid"
)

const (
	// defaultForbiddenHTTPMessage defines the generic message to not expose sensitive authorization information
	defaultForbiddenHTTPMessage = "request not authorized"
	// defaultInternalErrorMessage defines the generic message to return in case an error cannot be directly fixed by the user by updating their request
	defaultInternalErrorMessage = "internal error"
)

// signerErrorTypeToHTTPErrorStatus defines the correlation between signer dsl errors to HTTPErrorStatus
var signerErrorTypeToHTTPErrorStatus = map[signererrors.ErrorType]HTTPErrorStatus{
	signererrors.ErrInvalidArgument:    StatusInvalidArgument,
	signererrors.ErrNotFound:           StatusNotFound,
	signererrors.ErrAlreadyExists:      StatusAlreadyExists,
	signererrors.ErrInternal:           StatusInternal,
	signererrors.ErrNotImplemented:     StatusNotImplemented,
	signererrors.ErrBadGateway:         StatusBadGateway,
	signererrors.ErrPreconditionFailed: StatusPreconditionFailed,
	signererrors.ErrPermissionDenied:   StatusPermissionDenied,
}

// signerErrorTypeToHTTPErrorStatus defines the correlation between signer dsl errors to HTTPErrorStatus
var httpErrorStatusToHTTPCode = map[HTTPErrorStatus]int{
	StatusInvalidArgument:    http.StatusBadRequest,
	StatusAlreadyExists:      http.StatusBadRequest,
	StatusNotFound:           http.StatusNotFound,
	StatusInternal:           http.StatusInternalServerError,
	StatusNotImplemented:     http.StatusNotImplemented,
	StatusBadGateway:         http.StatusBadRequest,
	StatusPreconditionFailed: http.StatusPreconditionFailed,
	StatusPermissionDenied:   http.StatusForbidden,
}

// httpErrorDetails defines the details of an HTTPError
type httpErrorDetails struct {
	// message string message
	message string
	// traceableErrorID traceable error identifier
	traceableErrorID string
	// originalError defines the original error
	originalError error
}

// HTTPError defines a response for a failed HTTP call
type HTTPError struct {
	// code error code
	code int
	// status error status
	status string
	// details error details
	details httpErrorDetails
}

// HTTPErrorStatus defines the possible error statuses returned through the API. They match the ones defined in the API spec.
type HTTPErrorStatus string

// HTTPErrorStatus as defined in the API spec.
const (
	StatusInvalidArgument    HTTPErrorStatus = "INVALID_ARGUMENT"
	StatusAlreadyExists      HTTPErrorStatus = "ALREADY_EXISTS"
	StatusPermissionDenied   HTTPErrorStatus = "PERMISSION_DENIED"
	StatusNotFound           HTTPErrorStatus = "NOT_FOUND"
	StatusPreconditionFailed HTTPErrorStatus = "PRECONDITION_FAILED"
	StatusInternal           HTTPErrorStatus = "INTERNAL"
	StatusNotImplemented     HTTPErrorStatus = "NOT_IMPLEMENTED"
	StatusBadGateway         HTTPErrorStatus = "BAD_GATEWAY"
)

// NewHTTPError creates a new HTTPError given a code and a status
func NewHTTPError(status HTTPErrorStatus) *HTTPError {
	// this check is to be protected as there could be a mismatch between signerErrorTypeToHTTPErrorStatus and httpErrorStatusToHTTPCode
	code, ok := httpErrorStatusToHTTPCode[status]
	if !ok {
		status = StatusInternal
		code = http.StatusInternalServerError
	}

	return &HTTPError{
		code:   code,
		status: string(status),
	}
}

// NewHTTPErrorFromUseCaseError creates an HTTPError given an error returned from a use case.
func NewHTTPErrorFromUseCaseError(ctx context.Context, originalErr error) *HTTPError {
	useCaseError, ok := signererrors.CastAsUseCaseError(originalErr)
	if !ok {
		return NewHTTPErrorFromError(ctx, originalErr, StatusInternal)
	}

	status, ok := signerErrorTypeToHTTPErrorStatus[useCaseError.Type()]
	if !ok {
		status = StatusInternal
	}
	httpError := NewHTTPError(status).SetMessage(parseErrorMessage(useCaseError))
	httpError.SetOriginalError(useCaseError)
	return httpError
}

func parseErrorMessage(signerErr signererrors.UseCaseError) string {
	if signererrors.IsPermissionDenied(signerErr) {
		return defaultForbiddenHTTPMessage
	}
	if signererrors.IsInternal(signerErr) && signerErr.HumanReadableMessage() == nil {
		return defaultInternalErrorMessage
	}

	humanReadableMessage := signerErr.HumanReadableMessage()
	if humanReadableMessage != nil {
		return *humanReadableMessage
	}
	return ""
}

// NewHTTPErrorFromError creates a new error from a regular GO error interface
func NewHTTPErrorFromError(ctx context.Context, originalErr error, status HTTPErrorStatus) *HTTPError {
	httpError := NewHTTPError(status)
	httpError.SetOriginalError(originalErr)

	spanID, err := requestcontext.SpanIDFromContext(ctx)
	if err == nil {
		httpError.SetTraceableErrorID(*spanID)
	}

	return httpError
}

// Error returns the message of the error
func (e *HTTPError) Error() string {
	return e.details.message
}

// Code returns the HTTP status code corresponding to the error
func (e *HTTPError) Code() int {
	return e.code
}

// SetMessage sets the error message
func (e *HTTPError) SetMessage(msg string) *HTTPError {
	e.details.message = msg
	return e
}

// SetTraceableErrorID sets traceableErrrorID
func (e *HTTPError) SetTraceableErrorID(traceableErrorID string) {
	e.details.traceableErrorID = traceableErrorID
}

// SetOriginalError sets the original error of a wrapped error
func (e *HTTPError) SetOriginalError(err error) {
	e.details.originalError = err
}

// OriginalError returns the original error
func (e *HTTPError) OriginalError() error {
	return e.details.originalError
}

func (e *HTTPError) toErrorResponse(ctx context.Context) ErrorResponse {
	e.logError(ctx)
	return ErrorResponse{
		Code:   e.code,
		Status: e.status,
		Details: ErrorResponseDetails{
			TraceableErrorId: e.details.traceableErrorID,
			Message:          e.details.message,
		},
	}
}

func (e *HTTPError) logError(ctx context.Context) {
	errorID := uuid.New().String()
	e.SetTraceableErrorID(errorID)

	originalErr := e.details.originalError
	if originalErr == nil {
		return
	}

	logEntry := logger.LogEntry(ctx)
	logEntry.WithArguments("traceableErrorId", errorID)

	useCaseError, ok := signererrors.CastAsUseCaseError(originalErr)
	if !ok {
		if e.code == http.StatusInternalServerError {
			logEntry.Error(originalErr.Error())
			return
		}
		logEntry.Debug(originalErr.Error())
		return
	}

	if e.code == http.StatusInternalServerError {
		logEntry.Error(useCaseError.GetStack())
		return
	}
	logEntry.Debug(useCaseError.GetStack())
}
