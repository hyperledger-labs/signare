package errors

import (
	"fmt"
)

// PublicError is a traceable error that allows to edit the message to explain the reason behind it.
type PublicError struct {
	PrivateError
}

// WithMessage allows to extend the error message.
func (e *PublicError) WithMessage(format string, args ...any) *PublicError {
	e.errMessage = fmt.Sprintf(format, args...)
	return e
}

// SetHumanReadableMessage allows to set a message associated with the error that can be read from outside the application.
// Only UseCases should set human-readable messages.
func (e *PublicError) SetHumanReadableMessage(format string, args ...any) *PublicError {
	msg := fmt.Sprintf(format, args...)
	e.humanReadableMessage = &msg
	return e
}
