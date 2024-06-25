// Package referentialintegritydbout defines the output database adapters for the Referential Integrity.
package referentialintegritydbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/referentialintegritydb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

var _ referentialintegrity.ReferentialIntegrityStorage = new(Repository)

// Add a Referential Integrity Entry to the storage.
func (r *Repository) Add(ctx context.Context, data referentialintegrity.ReferentialIntegrityEntry) (*referentialintegrity.ReferentialIntegrityEntry, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	entryDB, addErr := r.infra.Add(ctx, *db)
	if addErr != nil {
		return nil, mapPersistenceErrorToSignerError(addErr)
	}

	addedHSM, err := mapFromDB(*entryDB)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedHSM, nil
}

// Get a Referential Integrity Entry from storage.
func (r *Repository) Get(ctx context.Context, id entities.StandardID) (*referentialintegrity.ReferentialIntegrityEntry, error) {
	storageData, err := r.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource referential integrity entry does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining referential integrity entry")
	}

	storedEntry, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedEntry, nil
}

// Remove removes a Referential Integrity Entry from storage.
func (r *Repository) Remove(ctx context.Context, id entities.StandardID) (*referentialintegrity.ReferentialIntegrityEntry, error) {
	storageData, err := r.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = r.infra.Remove(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return storageData, nil
}

// RemoveAllFromResource removes all entries where the resource  from storage.
func (r *Repository) RemoveAllFromResource(ctx context.Context, resourceID, resourceKind string) (*referentialintegrity.ReferentialIntegrityEntryCollection, error) {
	filter := r.Filter()
	filter.FilterByResource(resourceID, resourceKind)
	storageData, err := r.All(ctx, filter)
	if err != nil {
		return nil, err
	}

	_, err = r.infra.RemoveAllFromResource(ctx, resourceID, resourceKind)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return storageData, nil
}

// All Referential Integrity Entries from storage according to the provided filters.
func (r *Repository) All(ctx context.Context, filters referentialintegrity.ReferentialIntegrityFilters) (*referentialintegrity.ReferentialIntegrityEntryCollection, error) {
	f, ok := filters.(*referentialIntegrityEntryDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := r.infra.All(ctx, *f.ReferentialIntegrityEntryDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := referentialintegrity.ReferentialIntegrityEntryCollection{}
	if f.Pagination != nil {
		collection.Offset = f.Pagination.Offset
		collection.Limit = f.Pagination.Limit - 1
		if len(storageData) == f.Pagination.Limit {
			collection.MoreItems = true
			storageData = storageData[:len(storageData)-1]
		}
		f.Pagination.Limit--
	} else {
		collection.StandardCollectionPage = entities.NewUnlimitedQueryStandardCollectionPage(len(storageData))
	}

	items, err := mapSliceFromDB(storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	collection.Items = items

	return &collection, nil
}

func (r *Repository) Filter() referentialintegrity.ReferentialIntegrityFilters {
	storageFilter := referentialIntegrityEntryDBFilter{
		ReferentialIntegrityEntryDBFilter: &referentialintegritydb.ReferentialIntegrityEntryDBFilter{},
	}
	return &storageFilter
}

// Repository implementation of referentialintegrity.ReferentialIntegrityStorage
type Repository struct {
	infra *referentialintegritydb.ReferentialIntegrityEntryRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *referentialintegritydb.ReferentialIntegrityEntryRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ referentialintegrity.ReferentialIntegrityFilters = new(referentialIntegrityEntryDBFilter)

// FilterByParent filters the items by the specified parent resource ID and Kind.
func (filter *referentialIntegrityEntryDBFilter) FilterByParent(parentResourceID string, parentResourceKind string) referentialintegrity.ReferentialIntegrityFilters {
	filter.ParentResourceID = parentResourceID
	filter.ParentResourceKind = parentResourceKind
	filter.AppendFilter(postgres.NewEqualFilter("parent_resource_id"))
	filter.AppendFilter(postgres.NewEqualFilter("parent_resource_kind"))
	return filter
}

// FilterByResource filters the items by the specified resource ID and Kind.
func (filter *referentialIntegrityEntryDBFilter) FilterByResource(resourceID string, resourceKind string) referentialintegrity.ReferentialIntegrityFilters {
	filter.ResourceID = resourceID
	filter.ResourceKind = resourceKind
	filter.AppendFilter(postgres.NewEqualFilter("resource_id"))
	filter.AppendFilter(postgres.NewEqualFilter("resource_kind"))
	return filter
}

// OrderByCreationDate orders resources in storage by creation date.
func (filter *referentialIntegrityEntryDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) referentialintegrity.ReferentialIntegrityFilters {
	filter.ReferentialIntegrityEntryDBFilter = filter.ReferentialIntegrityEntryDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date.
func (filter *referentialIntegrityEntryDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) referentialintegrity.ReferentialIntegrityFilters {
	filter.ReferentialIntegrityEntryDBFilter = filter.ReferentialIntegrityEntryDBFilter.Sort("last_update", orderDirection)
	return filter
}

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
func (filter *referentialIntegrityEntryDBFilter) Paged(limit int, offset int) referentialintegrity.ReferentialIntegrityFilters {
	filter.ReferentialIntegrityEntryDBFilter = filter.ReferentialIntegrityEntryDBFilter.Paged(limit, offset)
	return filter
}

type referentialIntegrityEntryDBFilter struct {
	*referentialintegritydb.ReferentialIntegrityEntryDBFilter
}
