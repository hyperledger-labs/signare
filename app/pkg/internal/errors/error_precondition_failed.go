package errors

import (
	"errors"
	"fmt"
)

// PreconditionFailed returns a new PrivateError of type ErrPreconditionFailed.
func PreconditionFailed() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrPreconditionFailed),
			errType: ErrPreconditionFailed,
			stack:   getStack(),
		},
	}
}

// PreconditionFailedFromErr returns a new PrivateError of type ErrPreconditionFailed wrapping the original error.
func PreconditionFailedFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrPreconditionFailed),
			errType:    ErrPreconditionFailed,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsPreconditionFailed checks whether the target error is of type ErrPreconditionFailed.
func IsPreconditionFailed(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrPreconditionFailed
	}
	return false
}
