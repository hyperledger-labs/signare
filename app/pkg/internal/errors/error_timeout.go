package errors

import (
	"errors"
	"fmt"
)

// Timeout returns a new PrivateError of type ErrTimeout.
func Timeout() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrTimeout),
			errType: ErrTimeout,
			stack:   getStack(),
		},
	}
}

// TimeoutFromErr returns a new PrivateError of type ErrTimeout wrapping the original error.
func TimeoutFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrTimeout),
			errType:    ErrTimeout,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsTimeout checks whether the target error is of type ErrTimeout.
func IsTimeout(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrTimeout
	}
	return false
}
