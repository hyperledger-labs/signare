// Package admindbout defines the output database adapters for the Admin resource.
package admindbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/admindb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
)

var _ admin.AdminStorage = new(Repository)

// Add a Admin to storage.
func (repository *Repository) Add(ctx context.Context, data admin.Admin) (*admin.Admin, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	storageData, err := repository.infra.Add(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	addedAdmin, err := mapFromDB(*storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedAdmin, nil
}

// Get a Admin from storage.
func (repository *Repository) Get(ctx context.Context, id entities.StandardID) (*admin.Admin, error) {
	storageData, err := repository.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource 'admin' does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'admin'")
	}

	storedAdmin, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedAdmin, nil
}

// Edit a Admin from in storage.
func (repository *Repository) Edit(ctx context.Context, data admin.Admin) (*admin.Admin, error) {
	db, err := mapToUpdateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	result, err := repository.infra.Edit(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	rowsAffected, errRowsAffected := result.Result.RowsAffected()
	if errRowsAffected != nil {
		return nil, errors.InternalFromErr(err)
	}

	if rowsAffected == 0 {
		return nil, errors.NotFound().WithMessage("resource 'admin' does not match the one stored")
	}

	if rowsAffected > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'admin'")
	}

	storageData, err := repository.Get(ctx, data.StandardID)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storageData, nil
}

// Remove a Admin from the storage.
func (repository *Repository) Remove(ctx context.Context, id entities.StandardID) (*admin.Admin, error) {
	storageData, err := repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = repository.infra.Remove(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return storageData, nil
}

// All retrieves all Admins from the storage.
func (repository *Repository) All(ctx context.Context, filters admin.AdminFilters) (*admin.AdminCollection, error) {
	f, ok := filters.(*adminDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := repository.infra.List(ctx, *f.AdminDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := admin.AdminCollection{}
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

// Filter creates a new filter for the provided admin.
func (repository *Repository) Filter() admin.AdminFilters {
	storageFilter := adminDBFilter{
		AdminDBFilter: &admindb.AdminDBFilter{},
	}
	return &storageFilter
}

// Repository implementation of admin.AdminStorage
type Repository struct {
	infra *admindb.AdminRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *admindb.AdminRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ admin.AdminFilters = (*adminDBFilter)(nil)

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter
func (filter *adminDBFilter) Paged(limit int, offset int) admin.AdminFilters {
	filter.AdminDBFilter = filter.AdminDBFilter.Paged(limit, offset)
	return filter
}

// OrderByCreationDate orders resources in storage by creation date
func (filter *adminDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) admin.AdminFilters {
	filter.AdminDBFilter = filter.AdminDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date
func (filter *adminDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) admin.AdminFilters {
	filter.AdminDBFilter = filter.AdminDBFilter.Sort("last_update", orderDirection)
	return filter
}

type adminDBFilter struct {
	*admindb.AdminDBFilter
}
