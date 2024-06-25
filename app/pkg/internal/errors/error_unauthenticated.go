package errors

import (
	"errors"
	"fmt"
)

// Unauthenticated returns a new PrivateError of type ErrUnauthenticated.
func Unauthenticated() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrUnauthenticated),
			errType: ErrUnauthenticated,
			stack:   getStack(),
		},
	}
}

// UnauthenticatedFromErr returns a new PrivateError of type ErrUnauthenticated wrapping the original error.
func UnauthenticatedFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrUnauthenticated),
			errType:    ErrUnauthenticated,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsUnauthenticated checks whether the target error is of type ErrUnauthenticated.
func IsUnauthenticated(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrUnauthenticated
	}
	return false
}
