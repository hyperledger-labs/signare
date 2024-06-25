package errors

import (
	"bytes"
	"errors"
	"fmt"
)

// PrivateError is a traceable error that does not allow to edit the message to explain the reason behind it.
type PrivateError struct {
	err                  error
	errMessage           string
	errType              ErrorType
	humanReadableMessage *string
	wrappedErr           error
	stack                []trace
}

type trace struct {
	function string
	file     string
	line     int
}

// PrivateError implements the built-in error interface.
// It returns all available error information, including inner errors that are wrapped by this error.
func (e *PrivateError) Error() string {
	var errMsg = e.err.Error()
	if len(e.errMessage) > 0 {
		errMsg = fmt.Sprintf("%s: %s", errMsg, e.errMessage)
	}
	if e.wrappedErr == nil {
		return errMsg
	}
	return fmt.Sprintf("%s (wrapped error: %s)", errMsg, e.wrappedErr.Error())
}

// WithMessage allows to extend the error message.
func (e *PrivateError) WithMessage(format string, args ...any) *PrivateError {
	e.errMessage = fmt.Sprintf(format, args...)
	return e
}

// HumanReadableMessage retrieves the message associated with the error or nil if there isn't one.
// This message may be outputted outside the application.
func (e *PrivateError) HumanReadableMessage() *string {
	return e.humanReadableMessage
}

// Type returns the PrivateError Type
func (e *PrivateError) Type() ErrorType {
	return e.errType
}

// GetStack returns the PrivateError's stack trace.
func (e *PrivateError) GetStack() string {
	err := *e
	var originalSignerErr *PrivateError
	var stackTraceBottom = false

	for !stackTraceBottom {
		if errors.As(err.wrappedErr, &originalSignerErr) {
			err.wrappedErr = originalSignerErr.wrappedErr
			continue
		}
		stackTraceBottom = true
	}

	// As there were no wrapped error of type [PrivateError], we use the original error for printing the stack
	if originalSignerErr == nil {
		originalSignerErr = &err
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	_, _ = fmt.Fprintf(buf, "Error: %s\n", e.Error())
	_, _ = fmt.Fprintf(buf, "Original Error Stack Trace:\n")
	for _, frame := range originalSignerErr.stack {
		_, _ = fmt.Fprintf(buf, "\tat %s\n", frame.function)
		_, _ = fmt.Fprintf(buf, "\t\t%s:%d\n", frame.file, frame.line)
	}
	return buf.String()
}
