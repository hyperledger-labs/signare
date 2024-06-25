package signaturemanager

import (
	"errors"
	"fmt"
)

// Error for the persistence framework
type Error struct {
	description string
	err         error
}

var (
	errLibFailed           = errors.New("unable to initialize library")
	errInvalidSlot         = errors.New("invalid slot")
	errPinIncorrect        = errors.New("the pin is incorrect")
	errAlreadyInitialized  = errors.New("library is already initialized")
	errKeyGenerationFailed = errors.New("key generation failed")
	errInternal            = errors.New("internal error")
	errNotFound            = errors.New("not found")
	errInvalidArgument     = errors.New("invalid argument")
)

func (e *Error) Error() string {
	if len(e.description) == 0 {
		return e.err.Error()
	}
	return fmt.Sprintf("%s: %s", e.err.Error(), e.description)
}

func (e *Error) WithMessage(message string) *Error {
	e.description = message
	return e
}

func NewLibFailedError() *Error {
	return &Error{
		err: errLibFailed,
	}
}

func NewInvalidSlotError() *Error {
	return &Error{
		err: errInvalidSlot,
	}
}

func NewPinIncorrectError() *Error {
	return &Error{
		err: errPinIncorrect,
	}
}

func NewAlreadyInitializedError() *Error {
	return &Error{
		err: errAlreadyInitialized,
	}
}

func NewKeyGenerationError() *Error {
	return &Error{
		err: errKeyGenerationFailed,
	}
}

func NewInternalError() *Error {
	return &Error{
		err: errInternal,
	}
}

func NewNotFoundError() *Error {
	return &Error{
		err: errNotFound,
	}
}

func NewInvalidArgumentError() *Error {
	return &Error{
		err: errInvalidArgument,
	}
}

func IsLibFailedFailedError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errLibFailed)
	}
	return false
}

func IsInvalidSlotError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errInvalidSlot)
	}
	return false
}

func IsPinIncorrectError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errPinIncorrect)
	}
	return false
}

func IsAlreadyInitializedErr(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errAlreadyInitialized)
	}
	return false
}

func IsKeyGenerationError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errKeyGenerationFailed)
	}
	return false
}

func IsInternalError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errInternal)
	}
	return false
}

func IsNotFoundError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errNotFound)
	}
	return false
}

func IsInvalidArgumentError(err error) bool {
	var pkcsErr *Error
	if errors.As(err, &pkcsErr) {
		return errors.Is(pkcsErr.err, errInvalidArgument)
	}
	return false
}
