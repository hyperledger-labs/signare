package errors

import (
	"errors"
	"fmt"
)

// AlreadyExists returns a new PrivateError of type ErrAlreadyExists.
func AlreadyExists() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrAlreadyExists),
			errType: ErrAlreadyExists,
			stack:   getStack(),
		},
	}
}

// AlreadyExistsFromErr returns a new PrivateError of type ErrAlreadyExists wrapping the original error.
func AlreadyExistsFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrAlreadyExists),
			errType:    ErrAlreadyExists,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsAlreadyExists checks whether the target error is of type ErrAlreadyExists.
func IsAlreadyExists(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrAlreadyExists
	}
	return false
}
