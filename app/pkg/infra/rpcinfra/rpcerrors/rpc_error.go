package rpcerrors

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	// Error codes as defined in https://www.jsonrpc.org/specification

	ParseErrorCode              ErrorCode = -32700
	InvalidRequestErrorCode     ErrorCode = -32600
	MethodNotFoundErrorCode     ErrorCode = -32601
	InvalidParamsErrorCode      ErrorCode = -32602
	InternalErrorCode           ErrorCode = -32603
	UnauthorizedErrorCode       ErrorCode = -32099
	NotFoundErrorCode           ErrorCode = -32098
	PreconditionFailedErrorCode ErrorCode = -32097

	ParseErrorMsg              ErrorMsg = "Parse error"
	InvalidRequestErrorMsg     ErrorMsg = "Invalid request"
	MethodNotFoundErrorMsg     ErrorMsg = "Method not found"
	InvalidParamsErrorMsg      ErrorMsg = "Invalid params"
	InternalErrorMsg           ErrorMsg = "Internal error"
	UnauthorizedErrorMsg       ErrorMsg = "Unauthorized"
	NotFoundErrorMsg           ErrorMsg = "Not found"
	PreconditionFailedErrorMsg ErrorMsg = "Precondition failed"
)

// NewMethodNotFound creates a new method not found RPCError.
func NewMethodNotFound() *RPCError {
	return &RPCError{
		Code:    MethodNotFoundErrorCode,
		Message: MethodNotFoundErrorMsg,
	}
}

// NewMethodNotFoundFromErr creates a new method not found RPCError wrapping the original error.
func NewMethodNotFoundFromErr(err error) *RPCError {
	return &RPCError{
		Code:       MethodNotFoundErrorCode,
		Message:    MethodNotFoundErrorMsg,
		WrappedErr: err,
	}
}

// NewInvalidParams creates a new invalid params RPCError.
func NewInvalidParams() *RPCError {
	return &RPCError{
		Code:    InvalidParamsErrorCode,
		Message: InvalidParamsErrorMsg,
	}
}

// NewInvalidParamsFromErr creates a new invalid params RPCError wrapping the original error.
func NewInvalidParamsFromErr(err error) *RPCError {
	return &RPCError{
		Code:       InvalidParamsErrorCode,
		Message:    InvalidParamsErrorMsg,
		WrappedErr: err,
	}
}

// NewParse creates a new parse RPCError.
func NewParse() *RPCError {
	return &RPCError{
		Code:    ParseErrorCode,
		Message: ParseErrorMsg,
	}
}

// NewParseFromErr creates a new parse RPCError wrapping the original error.
func NewParseFromErr(err error) *RPCError {
	return &RPCError{
		Code:       ParseErrorCode,
		Message:    ParseErrorMsg,
		WrappedErr: err,
	}
}

// NewInvalidRequest creates a new invalid request RPCError.
func NewInvalidRequest() *RPCError {
	return &RPCError{
		Code:    InvalidRequestErrorCode,
		Message: InvalidRequestErrorMsg,
	}
}

// NewInvalidRequestFromErr creates a new invalid request RPCError wrapping the original error.
func NewInvalidRequestFromErr(err error) *RPCError {
	return &RPCError{
		Code:       InvalidRequestErrorCode,
		Message:    InvalidRequestErrorMsg,
		WrappedErr: err,
	}
}

// NewInternal creates a new internal server RPCError.
func NewInternal() *RPCError {
	return &RPCError{
		Code:    InternalErrorCode,
		Message: InternalErrorMsg,
	}
}

// NewInternalFromErr creates a new internal server RPCError wrapping the original error.
func NewInternalFromErr(err error) *RPCError {
	return &RPCError{
		Code:       InternalErrorCode,
		Message:    InternalErrorMsg,
		WrappedErr: err,
	}
}

// NewUnauthorized creates a new no authorized RPCError.
func NewUnauthorized() *RPCError {
	return &RPCError{
		Code:    UnauthorizedErrorCode,
		Message: UnauthorizedErrorMsg,
	}
}

// NewUnauthorizedFromErr creates a new no authorized RPCError wrapping the original error.
func NewUnauthorizedFromErr(err error) *RPCError {
	return &RPCError{
		Code:       UnauthorizedErrorCode,
		Message:    UnauthorizedErrorMsg,
		WrappedErr: err,
	}
}

// NewNotFound creates a new method not found RPCError.
func NewNotFound() *RPCError {
	return &RPCError{
		Code:    NotFoundErrorCode,
		Message: NotFoundErrorMsg,
	}
}

// NewNotFoundFromErr creates a new not found RPCError wrapping the original error.
func NewNotFoundFromErr(err error) *RPCError {
	return &RPCError{
		Code:       NotFoundErrorCode,
		Message:    NotFoundErrorMsg,
		WrappedErr: err,
	}
}

// NewPreconditionFailed creates a new method not found RPCError.
func NewPreconditionFailed() *RPCError {
	return &RPCError{
		Code:    PreconditionFailedErrorCode,
		Message: PreconditionFailedErrorMsg,
	}
}

// NewPreconditionFailedFromErr creates a new not found RPCError wrapping the original error.
func NewPreconditionFailedFromErr(err error) *RPCError {
	return &RPCError{
		Code:       PreconditionFailedErrorCode,
		Message:    PreconditionFailedErrorMsg,
		WrappedErr: err,
	}
}

// CastAsRPCError casts the provided error as an RPC error type
func CastAsRPCError(err error) (*RPCError, bool) {
	var castedErr *RPCError
	ok := errors.As(err, &castedErr)
	return castedErr, ok
}

// ErrorCode defines a json rpc 2.0 error code.
type ErrorCode int

// ErrorMsg defines a json rpc 2.0 error message.
type ErrorMsg string

// RPCError as defined in: https://www.jsonrpc.org/specification
type RPCError struct {
	// Code indicates the error type that occurred.
	Code ErrorCode `json:"code"`
	// Message provides a short description of the error.
	Message ErrorMsg `json:"message"`
	// WrappedErr contains the original error (if any).
	WrappedErr error `json:"-"`
}

// Marshal implements json.Marshaller.
func (rpcError *RPCError) Marshal() *json.RawMessage {
	result, err := json.Marshal(rpcError)
	if err != nil {
		return nil
	}
	raw := json.RawMessage(result)
	return &raw
}

// Error implements error interface
func (rpcError *RPCError) Error() string {
	msg := fmt.Sprintf("RPC %s error [code %v]", rpcError.Message, rpcError.Code)
	if rpcError.WrappedErr != nil {
		msg = fmt.Sprintf("%s: %v", msg, rpcError.WrappedErr)
	}
	return msg
}
