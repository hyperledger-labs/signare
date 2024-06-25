package errors

import (
	"errors"
	"fmt"
)

// BadGateway returns a new PrivateError of type ErrBadGateway.
func BadGateway() *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:     fmt.Errorf("%s", ErrBadGateway),
			errType: ErrBadGateway,
			stack:   getStack(),
		},
	}
}

// BadGatewayFromErr returns a new PrivateError of type ErrBadGateway wrapping the original error.
func BadGatewayFromErr(err error) *PublicError {
	return &PublicError{
		PrivateError: PrivateError{
			err:        fmt.Errorf("%s", ErrBadGateway),
			errType:    ErrBadGateway,
			wrappedErr: err,
			stack:      getStack(),
		},
	}
}

// IsBadGateway checks whether the target error is of type ErrBadGateway.
func IsBadGateway(err error) bool {
	var targetErr *PublicError
	if errors.As(err, &targetErr) {
		return targetErr.Type() == ErrBadGateway
	}
	return false
}
