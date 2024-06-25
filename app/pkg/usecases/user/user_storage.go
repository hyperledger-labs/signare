package user

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// UserStorage defines the functionality to interact with the User in storage.
type UserStorage interface {
	// Add a User in storage.
	Add(ctx context.Context, data User) (*User, error)
	// Get a User in storage.
	Get(ctx context.Context, id entities.ApplicationStandardID) (*User, error)
	// Edit a User in storage.
	Edit(ctx context.Context, data User) (*User, error)
	// Remove a User in storage.
	Remove(ctx context.Context, id entities.ApplicationStandardID) (*User, error)
	// All User in storage.
	All(ctx context.Context, filters UserFilters) (*UserCollection, error)

	// Filter by applicationID plus other optional filters.
	Filter(applicationID string) UserFilters
}

// UserFilters defines filter options for retrieving Users from storage.
type UserFilters interface {
	// OrderByCreationDate orders User in storage by creation date.
	OrderByCreationDate(direction persistence.OrderDirection) UserFilters
	// OrderByLastUpdateDate orders User in storage by last update date.
	OrderByLastUpdateDate(direction persistence.OrderDirection) UserFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
	Paged(limit int, offset int) UserFilters
}
