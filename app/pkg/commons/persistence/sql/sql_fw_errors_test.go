package sql

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"

	"github.com/stretchr/testify/assert"
)

func TestError_IsError(t *testing.T) {
	err := persistence.NewConfigCanNotBeLoadedError()
	assert.True(t, persistence.IsConfigCanNotBeLoaded(err))

	err = persistence.NewStatementCouldNotBePreparedError()
	assert.True(t, persistence.IsStatementCouldNotBePrepared(err))

	err = persistence.NewStatementExecutionFailedError()
	assert.True(t, persistence.IsStatementExecutionFailed(err))

	err = persistence.NewDBResponseCanNotBeProcessedError()
	assert.True(t, persistence.IsDBResponseCouldNotBeProcessed(err))

	err = persistence.NewCanNotBeginTransactionError()
	assert.True(t, persistence.IsCanNotBeginTransaction(err))

	err = persistence.NewCanNotRollbackTransactionError()
	assert.True(t, persistence.IsCanNotRollbackTransaction(err))

	err = persistence.NewCanNotCommitTransactionError()
	assert.True(t, persistence.IsCanNotCommitTransaction(err))

	err = persistence.NewAlreadyExistsError()
	assert.True(t, persistence.IsAlreadyExists(err))

	err = persistence.NewTxNotInContextError()
	assert.True(t, persistence.IsTxNotInContext(err))

	err = persistence.NewTransientConnectionError()
	assert.True(t, persistence.IsTransientConnectionError(err))

	err = persistence.NewPermanentConnectionError()
	assert.True(t, persistence.IsPermanentConnectionError(err))
}

func TestError_Description(t *testing.T) {
	t.Run("empty description", func(t *testing.T) {
		expectedError := "configuration can not be loaded"

		err := persistence.NewConfigCanNotBeLoadedError()
		assert.Equal(t, expectedError, err.Error())
	})

	t.Run("with description", func(t *testing.T) {
		expectedDescription := "test description"
		expectedError := "configuration can not be loaded: " + expectedDescription

		err := persistence.NewConfigCanNotBeLoadedError().WithMessage(expectedDescription)
		assert.Equal(t, expectedError, err.Error())
	})

}
