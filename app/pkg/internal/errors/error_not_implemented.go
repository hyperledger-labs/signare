package errors

import (
	"errors"
	"fmt"
)

// NotImplemented returns a new PrivateError of type ErrNotImplemented.
func NotImplemented() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrNotImplemented),
			errType: ErrNotImplemented,
			stack:   getStack(),
		},
	}
}

// NotImplementedFromErr returns a new PrivateError of type ErrNotImplemented wrapping the original error.
func NotImplementedFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrNotImplemented),
			errType:    ErrNotImplemented,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsNotImplemented checks whether the target error is of type ErrNotImplemented.
func IsNotImplemented(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrNotImplemented
	}
	return false
}
