// Package transactionaldbout defines the output database adapters for db transactional functions.
package transactionaldbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"
)

var _ transactionalmanager.TransactionalStorage = new(TransactionalRepository)

// TransactionalRepository repository for transactional operations
type TransactionalRepository struct {
	storage persistence.Storage
}

// TransactionalRepositoryOptions configures an TransactionalRepositoryOptions
type TransactionalRepositoryOptions struct {
	Storage persistence.Storage
}

// NewTransactionalRepository creates a TransactionalRepository with the given options
func NewTransactionalRepository(options TransactionalRepositoryOptions) *TransactionalRepository {
	return &TransactionalRepository{
		storage: options.Storage,
	}
}

// BeginTransaction starts a storage operation that is transactional. It returns a storage failure if it fails
func (repository *TransactionalRepository) BeginTransaction(ctx context.Context) (context.Context, error) {
	output, err := repository.storage.BeginTransaction(ctx)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	return output, nil
}

// RollbackTransaction rollbacks to previous state a storage operation that is transactional. It returns a storage failure if it fails
func (repository *TransactionalRepository) RollbackTransaction(ctx context.Context) error {
	err := repository.storage.RollbackTransaction(ctx)
	if err != nil {
		return errors.InternalFromErr(err)
	}
	return nil
}

// CommitTransaction confirms a storage operation that is transactional. It returns a storage failure if it fails
func (repository *TransactionalRepository) CommitTransaction(ctx context.Context) error {
	err := repository.storage.CommitTransaction(ctx)
	if err != nil {
		return errors.InternalFromErr(err)
	}
	return nil
}

// GetTransaction returns the transaction from context (in case of success).
func (repository *TransactionalRepository) GetTransaction(ctx context.Context) (any, error) {
	output, err := repository.storage.GetTransaction(ctx)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	return output, nil
}
