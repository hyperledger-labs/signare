// Package errors provides a way to create, manipulate and manage typed errors
// which contains stack trace information.
//
// All signare code inside the domain and application layers must use it.
package errors

import (
	"errors"
	"runtime"
	"sync"
)

// ErrorType associated with the PrivateError.
type ErrorType string

type UseCaseError interface {
	error
	Type() ErrorType
	GetStack() string
	HumanReadableMessage() *string
}

const (
	ErrInvalidArgument    ErrorType = "INVALID_ARGUMENT"
	ErrPreconditionFailed ErrorType = "PRECONDITION_FAILED"
	ErrUnauthenticated    ErrorType = "UNAUTHENTICATED"
	ErrPermissionDenied   ErrorType = "PERMISSION_DENIED"
	ErrNotFound           ErrorType = "NOT_FOUND"
	ErrAlreadyExists      ErrorType = "ALREADY_EXISTS"
	ErrNotImplemented     ErrorType = "NOT_IMPLEMENTED"
	ErrBadGateway         ErrorType = "BAD_GATEWAY"
	ErrUnavailable        ErrorType = "UNAVAILABLE"
	ErrTimeout            ErrorType = "TIMEOUT"
	ErrTooManyReq         ErrorType = "TOO_MANY_REQ"
	ErrInternal           ErrorType = "INTERNAL"
)

// CastAsUseCaseError casts the provided error as an internal error type
func CastAsUseCaseError(err error) (UseCaseError, bool) {
	var exportableError *PublicError
	isExportableError := errors.As(err, &exportableError)
	if isExportableError {
		return exportableError, isExportableError
	}

	var useCaseError *PrivateError
	isUseCaseError := errors.As(err, &useCaseError)
	if isUseCaseError {
		return useCaseError, isUseCaseError
	}

	return nil, false
}

func (et ErrorType) String() string {
	return string(et)
}

func getStack() []trace {
	var stack []trace
	var framesOnce sync.Once

	framesOnce.Do(func() {
		pcs := make([]uintptr, 32)
		npcs := runtime.Callers(6, pcs)
		stack = make([]trace, 0, npcs)
		callers := pcs[:npcs]
		ci := runtime.CallersFrames(callers)

		for {
			frame, more := ci.Next()
			stack = append(stack, trace{
				function: frame.Function,
				file:     frame.File,
				line:     frame.Line,
			})
			if !more || frame.Function == "main.main" {
				break
			}
		}
	})
	return stack
}
