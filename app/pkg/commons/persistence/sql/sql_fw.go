package sql

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"text/template"
	"time"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"

	"github.com/jmoiron/sqlx"
)

// The max idle connections for the database/sql handle. See database/sql docs.
const defaultMaxIdleConnections = 2

// The max open connections for the database/sql handle. See database/sql docs.
const defaultMaxOpenConnections = 100

// The max connection lifetime for the database/sql handle in milliseconds. See database/sql docs.
const defaultMaxConnectionLifetime = 0

var _ persistence.Storage = (*Fw)(nil)

// Fw persistence framework for Storage data:
// - Allows transactional executions
// - Translates specific database errors
// - Allow templated statements
type Fw struct {
	db              *sqlx.DB
	config          persistence.StorageConfig
	errorTranslator ErrorTranslator
}

type ContextKey string

const contextTransactionKey ContextKey = "adhgo_db_tx_key"

// SQLClientParameters defines parameters that will be passed to database/sql DB to configure client connections settings
type SQLClientParameters struct {
	MaxIdleConnections    *int
	MaxOpenConnections    *int
	MaxConnectionLifetime *int
}

// FwOptions options to configure the PersistenceFw
type FwOptions struct {
	Connection          Connection
	SQLClientParameters SQLClientParameters
}

// NewPersistenceFw creates a PersistenceFw instance with the given options
func NewPersistenceFw(opts FwOptions) (*Fw, error) {
	if opts.Connection == nil {
		return nil, errors.New("mandatory 'Connection' can not be nil")
	}
	persistenceFw := &Fw{
		config:          persistence.NewEmptyStorageConfig(),
		db:              opts.Connection.GetDB(),
		errorTranslator: opts.Connection.GetErrorTranslator(),
	}

	persistenceFw.SetSQLClientParameters(opts.SQLClientParameters)
	return persistenceFw, nil
}

func (p *Fw) SetSQLClientParameters(parameters SQLClientParameters) {
	maxIdleConnections := ptrIntDefaultValue(parameters.MaxIdleConnections, defaultMaxIdleConnections)
	maxOpenConnections := ptrIntDefaultValue(parameters.MaxOpenConnections, defaultMaxOpenConnections)
	maxConnectionLifetime := ptrIntDefaultValue(parameters.MaxConnectionLifetime, defaultMaxConnectionLifetime)

	p.db.SetMaxIdleConns(maxIdleConnections)
	p.db.SetMaxOpenConns(maxOpenConnections)
	p.db.SetConnMaxLifetime(time.Duration(maxConnectionLifetime) * time.Millisecond)
}

// BeginTransaction starts a Storage operation that is transactional. It returns an error if it fails
func (p *Fw) BeginTransaction(ctx context.Context) (context.Context, error) {
	txFromContext := ctx.Value(contextTransactionKey)

	_, ok := txFromContext.(*sqlx.Tx)
	if ok {
		return nil, persistence.NewCanNotBeginTransactionError().WithMessage("transaction already exists")
	}

	// Maybe add isolation level using BeginTxx instead of default tx
	tx, err := p.db.Beginx()

	if err != nil {
		return nil, persistence.NewCanNotBeginTransactionError().WithMessage(fmt.Sprintf("%v", err))
	}

	return context.WithValue(ctx, contextTransactionKey, tx), nil
}

func txFromContext(ctx context.Context) *sqlx.Tx {
	txFromContext := ctx.Value(contextTransactionKey)

	tx, ok := txFromContext.(*sqlx.Tx)
	if !ok {
		return nil
	}

	return tx
}

// RollbackTransaction rollbacks to previous state a Storage operation that is transactional. It returns an error if it fails
func (p *Fw) RollbackTransaction(ctx context.Context) error {
	txFromContext := ctx.Value(contextTransactionKey)

	tx, ok := txFromContext.(*sqlx.Tx)
	if !ok {
		return persistence.NewCanNotRollbackTransactionError().WithMessage("transaction not in context")
	}

	err := tx.Rollback()
	if err != nil {
		return persistence.NewCanNotRollbackTransactionError().WithMessage(fmt.Sprintf("%v", err))
	}

	return nil
}

// CommitTransaction confirms a Storage operation that is transactional. It returns an error if it fails
func (p *Fw) CommitTransaction(ctx context.Context) error {
	txFromContext := ctx.Value(contextTransactionKey)

	tx, ok := txFromContext.(*sqlx.Tx)
	if !ok {
		return persistence.NewCanNotCommitTransactionError().WithMessage("transaction not in context")
	}

	err := tx.Commit()
	if err != nil {
		return persistence.NewCanNotCommitTransactionError().WithMessage(fmt.Sprintf("%v", err))
	}

	return nil
}

// QueryAll obtains all element from statement with arguments provided. It returns an error if it fails
func (p *Fw) QueryAll(ctx context.Context, stmtID string, args interface{}, dst interface{}) error {
	// Templatize
	stmtContent, ok := p.config.GetStatement(stmtID)
	if !ok {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("statement with id [%s] not found", stmtID))
	}

	dstType := reflect.TypeOf(dst)
	dstTypeKind := dstType.Kind()

	if dstTypeKind != reflect.Ptr {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("destination variable of type [%s] must be a pointer", dstTypeKind))
	}

	dstPointerToType := dstType.Elem()
	dstPointerToTypeKind := dstPointerToType.Kind()
	if dstPointerToTypeKind != reflect.Array && dstPointerToTypeKind != reflect.Slice {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("destination variable of type [%s] must be an array or a slice", dstPointerToTypeKind))
	}

	sliceElementType := dstType.Elem().Elem()
	sliceValue := reflect.ValueOf(dst).Elem()
	templateParsed, err := template.New("").Parse(stmtContent.Content)
	if err != nil {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("failed parsing statement as go template. Statement id '%s' (%s) error: %v", stmtID, stmtContent.Content, err))
	}

	var buffer bytes.Buffer
	outWriter := bufio.NewWriter(&buffer)
	err = templateParsed.Execute(outWriter, args)
	if err != nil {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("failed executing go template. Statement id '%s' (%s) error: %v", stmtID, stmtContent.Content, err))
	}

	err = outWriter.Flush()
	if err != nil {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("executed template could not be flushed. Statement id '%s' error: %v", stmtID, err))
	}

	newStmt := buffer.String()
	var query *sqlx.NamedStmt

	tx := txFromContext(ctx)
	if tx != nil {
		query, err = tx.PrepareNamedContext(ctx, newStmt)
	} else {
		query, err = p.db.PrepareNamedContext(ctx, newStmt)
	}
	if err != nil {
		return persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("statement id '%s' (%s) error: %s", stmtID, newStmt, err))
	}
	defer query.Close()

	rows, err := query.QueryxContext(ctx, args)
	if err != nil {
		defaultError := persistence.NewStatementExecutionFailedError().WithMessage(fmt.Sprintf("statement id '%s' (%s) error: %s", stmtID, newStmt, err))
		return p.errorTranslator.TranslateError(ctx, err, defaultError)
	}

	defer rows.Close()

	elementPtr := reflect.New(sliceElementType)

	for rows.Next() {
		err = rows.StructScan(elementPtr.Interface())
		if err != nil {
			return persistence.NewDBResponseCanNotBeProcessedError().WithMessage(fmt.Sprintf("row could not be processed. Statement id '%s' (%s) error: %v", stmtID, newStmt, err))
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.Indirect(elementPtr)))
	}

	return nil
}

// ExecuteStmt executes statement with arguments provided. It returns an error if it fails
func (p *Fw) ExecuteStmt(ctx context.Context, stmtID string, args interface{}) error {
	_, err := p.ExecuteStmtWithStorageResult(ctx, stmtID, args)
	return err
}

// ExecuteStmtWithStorageResult executes statement with arguments provided. It returns information on stmt or an error if it fails
func (p *Fw) ExecuteStmtWithStorageResult(ctx context.Context, stmtID string, args interface{}) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	stmtContent, ok := p.config.GetStatement(stmtID)
	if !ok {
		return nil, persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("statement with id [%s] not found", stmtID))
	}
	templateParsed, err := template.New("").Parse(stmtContent.Content)
	if err != nil {
		return nil, persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("failed parsing statement as go template. Statement id '%s' (%s) error: %v", stmtID, stmtContent.Content, err))
	}

	var buffer bytes.Buffer
	outWriter := bufio.NewWriter(&buffer)
	err = templateParsed.Execute(outWriter, args)
	if err != nil {
		return nil, persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("failed executing go template. Statement id '%s' (%s) error: %v", stmtID, stmtContent.Content, err))
	}
	err = outWriter.Flush()

	if err != nil {
		return nil, persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("executed template could not be flushed. Statement id '%s' error: %v", stmtID, err))
	}

	newStmt := buffer.String()
	var namedStmt *sqlx.NamedStmt

	tx := txFromContext(ctx)
	if tx != nil {
		namedStmt, err = tx.PrepareNamedContext(ctx, newStmt)
	} else {
		namedStmt, err = p.db.PrepareNamedContext(ctx, newStmt)
	}
	if err != nil {
		return nil, persistence.NewStatementCouldNotBePreparedError().WithMessage(fmt.Sprintf("statement id '%s' (%s) error: %s", stmtID, newStmt, err))
	}
	defer namedStmt.Close()

	result, err := namedStmt.Exec(args)
	if err != nil {
		defaultError := persistence.NewStatementExecutionFailedError().WithMessage(fmt.Sprintf("statement id '%s' (%s) error: %s", stmtID, newStmt, err))
		return nil, p.errorTranslator.TranslateError(ctx, err, defaultError)
	}

	return &persistence.ExecuteStmtWithStorageResultOutput{Result: result}, nil
}

// GetTransaction returns the transaction from context (in case of success).
func (p *Fw) GetTransaction(ctx context.Context) (any, error) {
	tx := txFromContext(ctx)
	if tx == nil {
		return nil, persistence.NewTxNotInContextError()
	}
	return tx, nil

}

// AddConfig registers SQL statements
func (p *Fw) AddConfig(config persistence.StorageConfig) error {
	return p.config.AddConfig(config)
}

// PtrIntDefaultValue returns the default value if ptr is nil and the int value if not
func ptrIntDefaultValue(ptr *int, defaultValue int) int {
	if ptr != nil {
		return *ptr
	}
	return defaultValue
}
