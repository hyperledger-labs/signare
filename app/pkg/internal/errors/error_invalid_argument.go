package errors

import (
	"errors"
	"fmt"
)

// InvalidArgument returns a new PrivateError of type ErrInvalidArgument.
func InvalidArgument() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrInvalidArgument),
			errType: ErrInvalidArgument,
			stack:   getStack(),
		},
	}
}

// InvalidArgumentFromErr returns a new PrivateError of type ErrInvalidArgument wrapping the original error.
func InvalidArgumentFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrInvalidArgument),
			errType:    ErrInvalidArgument,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsInvalidArgument checks whether the target error is of type ErrInvalidArgument.
func IsInvalidArgument(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrInvalidArgument
	}
	return false
}
