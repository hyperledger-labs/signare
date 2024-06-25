package persistence

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
	errConfigCanNotBeLoaded        = errors.New("configuration can not be loaded")
	errStatementCouldNotBePrepared = errors.New("statement could not be prepared for execution")
	errStatementExecutionFailed    = errors.New("statement execution failed")
	errDBResponseCanNotBeProcessed = errors.New("can not process database response")
	errCanNotBeginTransaction      = errors.New("can not begin transaction")
	errCanNotRollbackTransaction   = errors.New("can not rollback transaction")
	errCanNotCommitTransaction     = errors.New("can not commit transaction")
	errAlreadyExists               = errors.New("already exists")
	errNotFound                    = errors.New("resource not found")
	errTxNotInContext              = errors.New("transaction not in context")
	errEntryNotAdded               = errors.New("entry not added")

	errConnectionIsPermanent = errors.New("connection cannot be executed")
	errConnectionIsTransient = errors.New("connection failed but could be executed in the future")
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

func NewTransientConnectionError() *Error {
	return &Error{
		err: errConnectionIsTransient,
	}
}

func NewPermanentConnectionError() *Error {
	return &Error{
		err: errConnectionIsPermanent,
	}
}

func NewConfigCanNotBeLoadedError() *Error {
	return &Error{
		err: errConfigCanNotBeLoaded,
	}
}

func NewStatementCouldNotBePreparedError() *Error {
	return &Error{
		err: errStatementCouldNotBePrepared,
	}
}

func NewStatementExecutionFailedError() *Error {
	return &Error{
		err: errStatementExecutionFailed,
	}
}

func NewDBResponseCanNotBeProcessedError() *Error {
	return &Error{
		err: errDBResponseCanNotBeProcessed,
	}
}

func NewCanNotBeginTransactionError() *Error {
	return &Error{
		err: errCanNotBeginTransaction,
	}
}

func NewCanNotRollbackTransactionError() *Error {
	return &Error{
		err: errCanNotRollbackTransaction,
	}
}

func NewCanNotCommitTransactionError() *Error {
	return &Error{
		err: errCanNotCommitTransaction,
	}
}

func NewAlreadyExistsError() *Error {
	return &Error{
		err: errAlreadyExists,
	}
}

func NewNotFoundError() *Error {
	return &Error{
		err: errNotFound,
	}
}

func NewTxNotInContextError() *Error {
	return &Error{
		err: errTxNotInContext,
	}
}

func NewEntryNotAddedError() *Error {
	return &Error{
		err: errEntryNotAdded,
	}
}

func IsConfigCanNotBeLoaded(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errConfigCanNotBeLoaded)
	}
	return false
}

func IsStatementCouldNotBePrepared(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errStatementCouldNotBePrepared)
	}
	return false
}

func IsStatementExecutionFailed(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errStatementExecutionFailed)
	}
	return false
}

func IsDBResponseCouldNotBeProcessed(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errDBResponseCanNotBeProcessed)
	}
	return false
}

func IsCanNotBeginTransaction(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errCanNotBeginTransaction)
	}
	return false
}

func IsCanNotRollbackTransaction(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errCanNotRollbackTransaction)
	}
	return false
}

func IsCanNotCommitTransaction(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errCanNotCommitTransaction)
	}
	return false
}

func IsAlreadyExists(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errAlreadyExists)
	}
	return false
}

func IsNotFound(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errNotFound)
	}
	return false
}

func IsTxNotInContext(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errTxNotInContext)
	}
	return false
}

func IsEntryNotAdded(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errEntryNotAdded)
	}
	return false
}

func IsPermanentConnectionError(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errConnectionIsPermanent)
	}
	return false
}

func IsTransientConnectionError(err error) bool {
	var persistenceErr *Error
	if errors.As(err, &persistenceErr) {
		return errors.Is(persistenceErr.err, errConnectionIsTransient)
	}
	return false
}
