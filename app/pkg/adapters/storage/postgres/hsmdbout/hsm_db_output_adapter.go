// Package hsmdbout defines the output database adapters for the HSMModule resource.
package hsmdbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmmoduledb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
)

var _ hsmmodule.HSMModuleStorage = new(Repository)

// Add an HSM to the database
func (repository *Repository) Add(ctx context.Context, data hsmmodule.HSMModule) (*hsmmodule.HSMModule, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	hsmDB, addErr := repository.infra.Add(ctx, *db)
	if addErr != nil {
		return nil, mapPersistenceErrorToSignerError(addErr)
	}

	addedHSM, err := mapFromDB(*hsmDB)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return addedHSM, nil
}

// All HSM from storage according to the provided filters
func (repository *Repository) All(ctx context.Context, filters hsmmodule.HSMModuleFilters) (*hsmmodule.HSMModulesCollection, error) {
	f, ok := filters.(*hardwareSecurityModuleDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := repository.infra.All(ctx, *f.HardwareSecurityModuleDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := hsmmodule.HSMModulesCollection{}
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

// Get an HSM from storages
func (repository *Repository) Get(ctx context.Context, id entities.StandardID) (*hsmmodule.HSMModule, error) {
	storageData, err := repository.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource HSM does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining HSM")
	}

	storedHSM, err := mapFromDB(storageData[0])
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storedHSM, nil
}

// Edit an HSM in storage
func (repository *Repository) Edit(ctx context.Context, data hsmmodule.HSMModule) (*hsmmodule.HSMModule, error) {
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
		return nil, errors.NotFound().WithMessage("resource HSM does not match the one stored")
	}

	if rowsAffected > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining HSM")
	}

	storageData, err := repository.Get(ctx, data.StandardID)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storageData, nil
}

// Remove removes an HSM from storage
func (repository *Repository) Remove(ctx context.Context, id entities.StandardID) (*hsmmodule.HSMModule, error) {
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

// Exists checks if the HSM is present in storage
func (repository *Repository) Exists(ctx context.Context, id entities.StandardID) error {
	_, err := repository.infra.Exists(ctx, id)
	if err != nil {
		return mapPersistenceErrorToSignerError(err)
	}

	return nil
}

func (repository *Repository) Filter() hsmmodule.HSMModuleFilters {
	storageFilter := hardwareSecurityModuleDBFilter{
		HardwareSecurityModuleDBFilter: &hsmmoduledb.HardwareSecurityModuleDBFilter{},
	}
	return &storageFilter
}

// Repository implementation of hsmmodule.HSMModuleStorage
type Repository struct {
	infra *hsmmoduledb.HardwareSecurityModuleRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *hsmmoduledb.HardwareSecurityModuleRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ hsmmodule.HSMModuleFilters = (*hardwareSecurityModuleDBFilter)(nil)

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter
func (filter *hardwareSecurityModuleDBFilter) Paged(limit int, offset int) hsmmodule.HSMModuleFilters {
	filter.HardwareSecurityModuleDBFilter = filter.HardwareSecurityModuleDBFilter.Paged(limit, offset)
	return filter
}

// OrderByCreationDate orders resources in storage by creation date
func (filter *hardwareSecurityModuleDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) hsmmodule.HSMModuleFilters {
	filter.HardwareSecurityModuleDBFilter = filter.HardwareSecurityModuleDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date
func (filter *hardwareSecurityModuleDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) hsmmodule.HSMModuleFilters {
	filter.HardwareSecurityModuleDBFilter = filter.HardwareSecurityModuleDBFilter.Sort("last_update", orderDirection)
	return filter
}

type hardwareSecurityModuleDBFilter struct {
	*hsmmoduledb.HardwareSecurityModuleDBFilter
}
