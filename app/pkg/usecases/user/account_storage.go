package user

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
)

// AccountStorage defines the functionality to interact with Account in storage.
type AccountStorage interface {
	// Add an Account in storage.
	Add(ctx context.Context, data Account) (*Account, error)
	// Get an Account from storage.
	Get(ctx context.Context, id AccountID) (*Account, error)
	// Remove an Account from storage.
	Remove(ctx context.Context, id AccountID) (*Account, error)
	// RemoveAllForAddress all Account for a given address and application from storage.
	RemoveAllForAddress(ctx context.Context, applicationID string, address address.Address) (*AccountCollection, error)
	// All Accounts in storage.
	All(ctx context.Context, filters AccountFilters) (*AccountCollection, error)

	// Filter by applicationID plus other optional filters.
	Filter(applicationID string) AccountFilters
}

// AccountFilters defines filter options for retrieving Accounts from storage.
type AccountFilters interface {
	// FilterByUserID filter accounts by user ID.
	FilterByUserID(userID string) AccountFilters
	// FilterByAddress filter accounts by address.
	FilterByAddress(address string) AccountFilters
}
