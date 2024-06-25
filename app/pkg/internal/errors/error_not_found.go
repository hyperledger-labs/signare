package errors

import (
	"errors"
	"fmt"
)

// NotFound returns a new PrivateError of type ErrNotFound.
func NotFound() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrNotFound),
			errType: ErrNotFound,
			stack:   getStack(),
		},
	}
}

// NotFoundFromErr returns a new PrivateError of type ErrNotFound wrapping the original error.
func NotFoundFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrNotFound),
			errType:    ErrNotFound,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsNotFound checks whether the target error is of type ErrNotFound.
func IsNotFound(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrNotFound
	}
	return false
}
