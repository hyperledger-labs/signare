package errors

import (
	"errors"
	"fmt"
)

// TooManyReq returns a new PrivateError of type ErrTooManyReq.
func TooManyReq() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrTooManyReq),
			errType: ErrTooManyReq,
			stack:   getStack(),
		},
	}
}

// TooManyReqFromErr returns a new PrivateError of type ErrTooManyReq wrapping the original error.
func TooManyReqFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrTooManyReq),
			errType:    ErrTooManyReq,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsTooManyReq checks whether the target error is of type ErrTooManyReq.
func IsTooManyReq(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrTooManyReq
	}
	return false
}
