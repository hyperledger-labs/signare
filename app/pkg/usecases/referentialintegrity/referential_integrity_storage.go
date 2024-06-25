package referentialintegrity

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// ReferentialIntegrityStorage defines the functionality to interact with the ReferentialIntegrityEntry in storage.
type ReferentialIntegrityStorage interface {
	// Add a ReferentialIntegrityEntry in storage.
	Add(ctx context.Context, data ReferentialIntegrityEntry) (*ReferentialIntegrityEntry, error)
	// Get a ReferentialIntegrityEntry from storage.
	Get(ctx context.Context, id entities.StandardID) (*ReferentialIntegrityEntry, error)
	// Remove a ReferentialIntegrityEntry in storage.
	Remove(ctx context.Context, id entities.StandardID) (*ReferentialIntegrityEntry, error)
	// RemoveAllFromResource removes all entries where the resource  from storage.
	RemoveAllFromResource(ctx context.Context, resourceID, resourceKind string) (*ReferentialIntegrityEntryCollection, error)
	// All ReferentialIntegrityEntry in storage.
	All(ctx context.Context, filters ReferentialIntegrityFilters) (*ReferentialIntegrityEntryCollection, error)
	// Filter all ReferentialIntegrityEntry by specific properties.
	Filter() ReferentialIntegrityFilters
}

// ReferentialIntegrityFilters defines filter options for retrieving ReferentialIntegrityEntry from storage.
type ReferentialIntegrityFilters interface {
	// FilterByParent filters entries by the specified parent ID and Kind.
	FilterByParent(parentResourceID string, parentResourceKind string) ReferentialIntegrityFilters
	// FilterByResource filters entries by the specified ID and Kind.
	FilterByResource(resourceID string, resourceKind string) ReferentialIntegrityFilters
	// OrderByCreationDate orders ReferentialIntegrityEntry in storage by creation date.
	OrderByCreationDate(orderDirection persistence.OrderDirection) ReferentialIntegrityFilters
	// OrderByLastUpdateDate orders ReferentialIntegrityEntry in storage by last update date.
	OrderByLastUpdateDate(orderDirection persistence.OrderDirection) ReferentialIntegrityFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
	Paged(limit int, offset int) ReferentialIntegrityFilters
}
