package application

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// ApplicationStorage defines the functionality to interact with the Application in storage.
type ApplicationStorage interface {
	// Add an Application in storage. It returns an error if it fails
	Add(ctx context.Context, data Application) (*Application, error)
	// All Application in storage. It returns an error if it fails
	All(ctx context.Context, filter ApplicationFilters) (*ApplicationCollection, error)
	// Get an Application in storage. It returns an error if it fails
	Get(ctx context.Context, id entities.StandardID) (*Application, error)
	// Edit an Application in storage. It returns an error if it fails
	Edit(ctx context.Context, data Application) (*Application, error)
	// Remove an Application in storage. It returns an error if it fails
	Remove(ctx context.Context, id entities.StandardID) (*Application, error)
	// Exists returns whether the specified Application exists in storage. It returns an error if the operation fails.
	Exists(ctx context.Context, id entities.StandardID) error
	// Filter provides a DSL to filter, order and paginate All operation
	Filter() ApplicationFilters
}

// ApplicationFilters defines all possible options to list Application from storage
type ApplicationFilters interface {
	// OrderByCreationDate orders Application in storage by creation date
	OrderByCreationDate(orderDirection persistence.OrderDirection) ApplicationFilters
	// OrderByLastUpdateDate orders Application in storage by last update date
	OrderByLastUpdateDate(orderDirection persistence.OrderDirection) ApplicationFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter
	Paged(limit int, offset int) ApplicationFilters
}
