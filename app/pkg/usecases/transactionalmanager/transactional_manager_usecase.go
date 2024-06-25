package transactionalmanager

import (
	"context"
	"sync"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

// FuncToExecuteAfterCommitType callback function to execute after a database transaction commit
type FuncToExecuteAfterCommitType func(context.Context) error

// TransactionalManagerUseCase defines the functionality to execute a transaction in a transactional manner.
type TransactionalManagerUseCase interface {
	// ExecuteInTransaction executes a database transaction.
	ExecuteInTransaction(ctx context.Context, transactionalFunc func(context.Context) (interface{}, error)) (interface{}, error)
	// RegisterAfterCommitIfTransactionInProgress callback to execute after committing a database transaction.
	RegisterAfterCommitIfTransactionInProgress(ctx context.Context, funcToExecuteAfterCommit FuncToExecuteAfterCommitType) bool
}

// ExecuteInTransaction begins the execution of a transaction and rollbacks (reverts) any persistent changes if the transaction fails.
// ctx determines the context to execute transaction within and transactionalFunc the transaction to perform.
// It returns transaction result as interface{} and it returns an usecase failure if fails.
func (t *TransactionalManager) ExecuteInTransaction(ctx context.Context, transactionalFunc func(context.Context) (interface{}, error)) (response interface{}, failure error) {
	txCtx := ctx
	_, err := t.transactionalStorage.GetTransaction(ctx)
	ongoingTransaction := true // if there are nested transactions, the first one is responsible for executing the final commit or rollback. We handle that with this variable
	if err != nil {            // tx is nil
		ongoingTransaction = false
		beginTransactionCtx, errBeginTransaction := t.transactionalStorage.BeginTransaction(ctx)
		if errBeginTransaction != nil {
			return nil, errBeginTransaction
		}
		txCtx = beginTransactionCtx
	}
	tx, err := t.transactionalStorage.GetTransaction(txCtx)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	defer func() {
		if x := recover(); x != nil {
			if !ongoingTransaction {
				// if a panic occurs, we try to roll the transaction back
				rollbackTransactionErr := t.transactionalStorage.RollbackTransaction(txCtx)
				failure = errors.Internal().WithMessage("Database transaction was rolled back")
				if rollbackTransactionErr != nil {
					// If the transaction has already been committed, we will get an error when trying to roll it back.
					// In this situation the panic occurred when executing the function after the commit.
					failure = rollbackTransactionErr
				}
			}
		}
	}()

	response, failure = transactionalFunc(txCtx)
	if failure != nil {
		if ongoingTransaction {
			return response, failure
		}
		err2 := t.transactionalStorage.RollbackTransaction(txCtx)
		if err2 != nil {
			logger.LogEntry(ctx).Error(err2.Error())
		}
		return nil, failure
	}

	if !ongoingTransaction {
		commitTransactionErr := t.transactionalStorage.CommitTransaction(txCtx)
		if commitTransactionErr != nil {
			return nil, commitTransactionErr
		}

		value, needToExecuteFunc := funcToExecuteAfterCommitMap.LoadAndDelete(tx)
		if needToExecuteFunc {
			funcToExecuteAfterCommitArray, correctValueType := value.([]FuncToExecuteAfterCommitType)
			if correctValueType {
				for _, funcToExecute := range funcToExecuteAfterCommitArray {
					failureFunc := funcToExecute(txCtx)
					if failureFunc != nil {
						return nil, failureFunc
					}
				}
			}
		}
	}

	return response, nil
}

// RegisterAfterCommitIfTransactionInProgress registers a callback function to execute after the database transaction in the provided context is committed
func (t *TransactionalManager) RegisterAfterCommitIfTransactionInProgress(ctx context.Context, funcToExecuteAfterCommit FuncToExecuteAfterCommitType) bool {
	tx, _ := t.transactionalStorage.GetTransaction(ctx)
	if tx == nil {
		// if transaction is nil, no transaction is in progress and no function can be registered
		return false
	}

	value, alreadyLoaded := funcToExecuteAfterCommitMap.Load(tx)
	if alreadyLoaded {
		funcToExecuteAfterCommitArray, correctValueType := value.([]FuncToExecuteAfterCommitType)
		if !correctValueType {
			logger.LogEntry(ctx).Warnf("Functions to execute after commit incorrectly registered [%s]", value)
			initialFuncToExecuteAfterCommitArray := []FuncToExecuteAfterCommitType{funcToExecuteAfterCommit}
			funcToExecuteAfterCommitMap.Store(tx, initialFuncToExecuteAfterCommitArray)
		} else {
			funcToExecuteAfterCommitMap.Store(tx, append(funcToExecuteAfterCommitArray, funcToExecuteAfterCommit))
		}
	} else {
		initialFuncToExecuteAfterCommitArray := []FuncToExecuteAfterCommitType{funcToExecuteAfterCommit}
		funcToExecuteAfterCommitMap.Store(tx, initialFuncToExecuteAfterCommitArray)
	}
	return true
}

var _ TransactionalManagerUseCase = new(TransactionalManager)

// TransactionalManager handles transactional functionality.
type TransactionalManager struct {
	transactionalStorage TransactionalStorage
}

// TransactionalManagerOptions configures a TransactionalManager.
type TransactionalManagerOptions struct {
	// TransactionalStorage transactional manager storage
	TransactionalStorage TransactionalStorage
}

// ProvideTransactionalManager creates a TransactionalManager object with given options.
func ProvideTransactionalManager(options TransactionalManagerOptions) *TransactionalManager {
	return &TransactionalManager{
		transactionalStorage: options.TransactionalStorage,
	}
}

var funcToExecuteAfterCommitMap sync.Map
