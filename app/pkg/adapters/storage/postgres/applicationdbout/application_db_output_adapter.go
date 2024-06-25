// Package applicationdbout defines the output database adapters for the Application resource.
package applicationdbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/applicationdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
)

var _ application.ApplicationStorage = new(Repository)

// Add an Application to the database
func (repository *Repository) Add(ctx context.Context, data application.Application) (*application.Application, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	storageData, err := repository.infra.Add(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	addedApplication, err := mapFromDB(*storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedApplication, nil
}

// All Application from storage according to the provided filters
func (repository *Repository) All(ctx context.Context, filters application.ApplicationFilters) (*application.ApplicationCollection, error) {
	f, ok := filters.(*applicationDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := repository.infra.All(ctx, *f.ApplicationDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := application.ApplicationCollection{}
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

// Get an Application from storages
func (repository *Repository) Get(ctx context.Context, id entities.StandardID) (*application.Application, error) {
	storageData, err := repository.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource application does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining application")
	}

	storedApplication, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedApplication, nil
}

// Edit an Application in storage
func (repository *Repository) Edit(ctx context.Context, data application.Application) (*application.Application, error) {
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
		return nil, errors.NotFound().WithMessage("resource application does not match the one stored")
	}

	if rowsAffected > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining application")
	}

	storageData, err := repository.Get(ctx, data.StandardID)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storageData, nil
}

// Exists returns whether the specified Application exists in storage. It returns an error if the operation fails.
func (repository *Repository) Exists(ctx context.Context, id entities.StandardID) error {
	_, err := repository.infra.Exists(ctx, id)
	if err != nil {
		signerErr := mapPersistenceErrorToSignerError(err)
		return signerErr
	}
	return nil
}

// Remove removes an Application from storage
func (repository *Repository) Remove(ctx context.Context, id entities.StandardID) (*application.Application, error) {
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

func (repository *Repository) Filter() application.ApplicationFilters {
	storageFilter := applicationDBFilter{
		ApplicationDBFilter: &applicationdb.ApplicationDBFilter{},
	}
	return &storageFilter
}

// Repository implementation of application.ApplicationStorage
type Repository struct {
	infra *applicationdb.ApplicationRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *applicationdb.ApplicationRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ application.ApplicationFilters = (*applicationDBFilter)(nil)

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter
func (filter *applicationDBFilter) Paged(limit int, offset int) application.ApplicationFilters {
	filter.ApplicationDBFilter = filter.ApplicationDBFilter.Paged(limit, offset)
	return filter
}

// OrderByCreationDate orders resources in storage by creation date
func (filter *applicationDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) application.ApplicationFilters {
	filter.ApplicationDBFilter = filter.ApplicationDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date
func (filter *applicationDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) application.ApplicationFilters {
	filter.ApplicationDBFilter = filter.ApplicationDBFilter.Sort("last_update", orderDirection)
	return filter
}

type applicationDBFilter struct {
	*applicationdb.ApplicationDBFilter
}
