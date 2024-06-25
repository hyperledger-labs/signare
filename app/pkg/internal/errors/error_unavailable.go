package errors

import (
	"errors"
	"fmt"
)

// Unavailable returns a new PrivateError of type ErrUnavailable.
func Unavailable() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrUnavailable),
			errType: ErrUnavailable,
			stack:   getStack(),
		},
	}
}

// UnavailableFromErr returns a new PrivateError of type ErrUnavailable wrapping the original error.
func UnavailableFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrUnavailable),
			errType:    ErrUnavailable,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsUnavailable checks whether the target error is of type ErrUnavailable.
func IsUnavailable(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrUnavailable
	}
	return false
}
