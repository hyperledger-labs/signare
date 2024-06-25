package admin

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// AdminStorage defines the functionality to interact with Admin in Storage.
type AdminStorage interface {
	// Add an Admin in AdminStorage.
	Add(ctx context.Context, data Admin) (*Admin, error)
	// Get an Admin from AdminStorage.
	Get(ctx context.Context, id entities.StandardID) (*Admin, error)
	// Edit and Admin from AdminStorage.
	Edit(ctx context.Context, data Admin) (*Admin, error)
	// Remove an Admin from AdminStorage.
	Remove(ctx context.Context, id entities.StandardID) (*Admin, error)
	// All Admins in AdminStorage.
	All(ctx context.Context, filters AdminFilters) (*AdminCollection, error)

	// Filter create an AdminFilters instance.
	Filter() AdminFilters
}

// AdminFilters defines filter options for retrieving Admins from Storage.
type AdminFilters interface {
	// OrderByCreationDate orders Admin in storage by creation date
	OrderByCreationDate(orderDirection persistence.OrderDirection) AdminFilters
	// OrderByLastUpdateDate orders Admin in storage by last update date
	OrderByLastUpdateDate(orderDirection persistence.OrderDirection) AdminFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter
	Paged(limit int, offset int) AdminFilters
}
