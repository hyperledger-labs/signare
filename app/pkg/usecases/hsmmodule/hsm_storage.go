package hsmmodule

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// HSMModuleStorage defines the functionality to interact with the CreateHSMModuleOutput in storage.
type HSMModuleStorage interface {
	// Add an HSMModule in storage. It returns an error if it fails.
	Add(ctx context.Context, data HSMModule) (*HSMModule, error)
	// All HSMModule in storage. It returns an error if it fails.
	All(ctx context.Context, filter HSMModuleFilters) (*HSMModulesCollection, error)
	// Get an HSMModule in storage. It returns an error if it fails.
	Get(ctx context.Context, id entities.StandardID) (*HSMModule, error)
	// Edit an HSMModule in storage. It returns an error if it fails.
	Edit(ctx context.Context, data HSMModule) (*HSMModule, error)
	// Remove an HSMModule in storage. It returns an error if it fails.
	Remove(ctx context.Context, id entities.StandardID) (*HSMModule, error)
	// Exists checks if an HSMModule is present in storage, returning an error if it is not.
	Exists(ctx context.Context, id entities.StandardID) error

	// Filter provides a DSL to filter, order and paginate All operation.
	Filter() HSMModuleFilters
}

// HSMModuleFilters defines all possible options to list CreateHSMModuleOutput from storage.
type HSMModuleFilters interface {
	// OrderByCreationDate orders CreateHSMModuleOutput in storage by creation date.
	OrderByCreationDate(orderDirection persistence.OrderDirection) HSMModuleFilters
	// OrderByLastUpdateDate orders CreateHSMModuleOutput in storage by last update date.
	OrderByLastUpdateDate(orderDirection persistence.OrderDirection) HSMModuleFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
	Paged(limit int, offset int) HSMModuleFilters
}
