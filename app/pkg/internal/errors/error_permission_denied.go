package errors

import (
	"errors"
	"fmt"
)

var defaultForbiddenHTTPMessage = "request not authorized"

// PermissionDenied returns a new PrivateError of type ErrPermissionDenied.
func PermissionDenied() *PrivateError {
	return &PrivateError{
		err:                  fmt.Errorf("%s", ErrPermissionDenied),
		errType:              ErrPermissionDenied,
		stack:                getStack(),
		humanReadableMessage: &defaultForbiddenHTTPMessage,
	}
}

// PermissionDeniedFromErr returns a new PrivateError of type ErrPermissionDenied wrapping the original error.
func PermissionDeniedFromErr(err error) *PrivateError {
	return &PrivateError{
		err:                  fmt.Errorf("%s", ErrPermissionDenied),
		errType:              ErrPermissionDenied,
		wrappedErr:           err,
		stack:                getStack(),
		humanReadableMessage: &defaultForbiddenHTTPMessage,
	}
}

// IsPermissionDenied checks whether the target error is of type ErrPermissionDenied.
func IsPermissionDenied(err error) bool {
	var targetErr *PrivateError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrPermissionDenied
	}
	return false
}
