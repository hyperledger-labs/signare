// Package userdbout defines the output database adapters for the User resource.
package userdbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/userdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

var _ user.UserStorage = new(Repository)

// Add a User to storage.
func (repository *Repository) Add(ctx context.Context, data user.User) (*user.User, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	storageData, err := repository.infra.Add(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	addedUser, err := mapFromDB(*storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedUser, nil
}

// Get a User from storage.
func (repository *Repository) Get(ctx context.Context, id entities.ApplicationStandardID) (*user.User, error) {
	storageData, err := repository.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource 'user' does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'user'")
	}

	storedUser, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedUser, nil
}

// Edit a User in storage.
func (repository *Repository) Edit(ctx context.Context, data user.User) (*user.User, error) {
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
		return nil, errors.NotFound().WithMessage("resource 'user' does not match the one stored")
	}

	if rowsAffected > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'user'")
	}

	storageData, err := repository.Get(ctx, data.ApplicationStandardID)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storageData, nil
}

// Remove a User from the storage.
func (repository *Repository) Remove(ctx context.Context, id entities.ApplicationStandardID) (*user.User, error) {
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

// All retrieves all Users from the storage.
func (repository *Repository) All(ctx context.Context, filters user.UserFilters) (*user.UserCollection, error) {
	f, ok := filters.(*userDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := repository.infra.List(ctx, *f.UserDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := user.UserCollection{}
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

// Filter creates a new filter for the provided application.
func (repository *Repository) Filter(applicationID string) user.UserFilters {
	storageFilter := userDBFilter{
		UserDBFilter: &userdb.UserDBFilter{
			UserDB: userdb.UserDB{
				ApplicationStandardID: entities.ApplicationStandardID{
					ApplicationID: applicationID,
				},
			},
		},
	}
	return &storageFilter
}

// Repository implementation of user.UserStorage
type Repository struct {
	infra *userdb.UserRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *userdb.UserRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ user.UserFilters = (*userDBFilter)(nil)

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
func (filter *userDBFilter) Paged(limit int, offset int) user.UserFilters {
	filter.UserDBFilter = filter.UserDBFilter.Paged(limit, offset)
	return filter
}

// OrderByCreationDate orders resources in storage by creation date.
func (filter *userDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) user.UserFilters {
	filter.UserDBFilter = filter.UserDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date.
func (filter *userDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) user.UserFilters {
	filter.UserDBFilter = filter.UserDBFilter.Sort("last_update", orderDirection)
	return filter
}

type userDBFilter struct {
	*userdb.UserDBFilter
}
