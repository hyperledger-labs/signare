package transactionalmanager

import (
	"context"
)

// TransactionalStorage defines a storage for transactional manager
type TransactionalStorage interface {
	// BeginTransaction starts a storage operation that is transactional. It returns an error if it fails
	BeginTransaction(ctx context.Context) (context.Context, error)
	// RollbackTransaction rollback to previous state a storage operation that is transactional. It returns an error if it fails
	RollbackTransaction(ctx context.Context) error
	// CommitTransaction confirms a storage operation that is transactional. It returns an error if it fails
	CommitTransaction(ctx context.Context) error
	// GetTransaction returns the transaction from context (in case of success).
	GetTransaction(ctx context.Context) (any, error)
}
