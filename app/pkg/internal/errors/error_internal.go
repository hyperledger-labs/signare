package errors

import (
	"errors"
	"fmt"
)

var defaultInternalErrorMessage = "internal error"

// Internal returns a new PrivateError of type ErrInternal.
func Internal() *PrivateError {
	return &PrivateError{
		err:                  fmt.Errorf("%s", ErrInternal),
		errType:              ErrInternal,
		stack:                getStack(),
		humanReadableMessage: &defaultInternalErrorMessage,
	}
}

// InternalFromErr returns a new PrivateError of type ErrInternal wrapping the original error.
func InternalFromErr(err error) *PrivateError {
	return &PrivateError{
		err:                  fmt.Errorf("%s", ErrInternal),
		errType:              ErrInternal,
		wrappedErr:           err,
		stack:                getStack(),
		humanReadableMessage: &defaultInternalErrorMessage,
	}
}

// IsInternal checks whether the target error is of type ErrInternal.
func IsInternal(err error) bool {
	var targetErr *PrivateError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrInternal
	}
	return false
}
