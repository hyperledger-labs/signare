// Package accountdbout defines the output database adapters for the Account resource.
package accountdbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/accountdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

var _ user.AccountStorage = new(Repository)

// Add an Account to the storage.
func (repository *Repository) Add(ctx context.Context, data user.Account) (*user.Account, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	storageData, err := repository.infra.Add(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	addedAccount, err := mapFromDB(*storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedAccount, nil
}

// Get an Account from the storage.
func (repository *Repository) Get(ctx context.Context, id user.AccountID) (*user.Account, error) {
	getInput := mapToAccountID(id)
	storageData, err := repository.infra.Get(ctx, getInput)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource 'account' does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining account")
	}

	storedAccount, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedAccount, nil
}

// All retrieves all Accounts from the storage.
func (repository *Repository) All(ctx context.Context, filters user.AccountFilters) (*user.AccountCollection, error) {
	f, ok := filters.(*accountDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	storageData, err := repository.infra.List(ctx, *f.AccountDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := user.AccountCollection{}
	collection.StandardCollectionPage = entities.NewUnlimitedQueryStandardCollectionPage(len(storageData))

	items, err := mapSliceFromDB(storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	collection.Items = items

	return &collection, nil
}

// Remove an Account from the storage.
func (repository *Repository) Remove(ctx context.Context, id user.AccountID) (*user.Account, error) {
	storageData, err := repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	removeInput := mapToAccountID(id)
	_, errRemove := repository.infra.Remove(ctx, removeInput)
	if errRemove != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return storageData, nil
}

// RemoveAllForAddress all Account for a given address and application from storage.
func (repository *Repository) RemoveAllForAddress(ctx context.Context, applicationID string, address address.Address) (*user.AccountCollection, error) {
	filters := repository.Filter(applicationID)
	filters.FilterByAddress(address.String())
	accounts, err := repository.All(ctx, filters)
	if err != nil {
		return nil, err
	}

	removeAllInput := mapToAccountApplicationAddresses(applicationID, address)
	_, errRemove := repository.infra.RemoveAllForAddress(ctx, removeAllInput)
	if errRemove != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return accounts, nil
}

// Filter creates a new filter for the provided application.
func (repository *Repository) Filter(applicationID string) user.AccountFilters {
	storageFilter := accountDBFilter{
		AccountDBFilter: &accountdb.AccountDBFilter{
			AccountDB: accountdb.AccountDB{
				ApplicationID: applicationID,
			},
		},
	}
	return &storageFilter
}

// Repository implementation of account.AccountStorage
type Repository struct {
	infra *accountdb.AccountRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *accountdb.AccountRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ user.AccountFilters = (*accountDBFilter)(nil)

func (filter *accountDBFilter) FilterByUserID(userID string) user.AccountFilters {
	filter.UserID = userID
	filter.AppendFilter(postgres.NewEqualFilter("user_id"))
	return filter
}

func (filter *accountDBFilter) FilterByAddress(address string) user.AccountFilters {
	filter.Address = address
	filter.AppendFilter(postgres.NewEqualFilter("address"))
	return filter
}

type accountDBFilter struct {
	*accountdb.AccountDBFilter
}
