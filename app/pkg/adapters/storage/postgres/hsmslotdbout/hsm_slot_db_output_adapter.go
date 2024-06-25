// Package hsmslotdbout defines the output database adapters for the HSMSlot resource.
package hsmslotdbout

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmslotdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
)

var _ hsmslot.HSMSlotStorage = new(Repository)

// Add an HSMSlot to storage.
func (r *Repository) Add(ctx context.Context, data hsmslot.HSMSlot) (*hsmslot.HSMSlot, error) {
	db, err := mapToCreateDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	storageData, err := r.infra.Add(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	return mapFromDB(*storageData)
}

// Get an HSMSlot from storage.
func (r *Repository) Get(ctx context.Context, id entities.StandardID) (*hsmslot.HSMSlot, error) {
	storageData, err := r.infra.Get(ctx, id)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource 'hsm_slot' does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'hsm_slot'")
	}

	return mapFromDB(storageData[0])
}

func (r *Repository) GetByApplication(ctx context.Context, applicationID entities.StandardID) (*hsmslot.HSMSlot, error) {
	storageData, err := r.infra.GetByApplicationID(ctx, applicationID)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	if len(storageData) == 0 {
		return nil, errors.NotFound().WithMessage("resource 'hsm_slot' does not exist")
	}

	if len(storageData) > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'hsm_slot'")
	}

	return mapFromDB(storageData[0])
}

// EditPin edits an HSMSlot's Pin in storage.
func (r *Repository) EditPin(ctx context.Context, data hsmslot.HSMSlot) (*hsmslot.HSMSlot, error) {
	db, err := mapToUpdatePinDB(data)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	result, err := r.infra.EditPin(ctx, *db)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	rowsAffected, errRowsAffected := result.Result.RowsAffected()
	if errRowsAffected != nil {
		return nil, errors.InternalFromErr(err)
	}

	if rowsAffected == 0 {
		return nil, errors.NotFound().WithMessage("resource 'hsm_slot' does not match the one stored")
	}

	if rowsAffected > 1 {
		return nil, errors.Internal().WithMessage("unexpected number of results when obtaining 'hsm_slot'")
	}

	storageData, err := r.Get(ctx, entities.StandardID{ID: data.ID})
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return storageData, nil
}

// Remove an HSMSlot from the storage.
func (r *Repository) Remove(ctx context.Context, id entities.StandardID) (*hsmslot.HSMSlot, error) {
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

// All retrieves all HSMSlots from the storage.
func (r *Repository) All(ctx context.Context, filters hsmslot.HSMSlotFilters) (*hsmslot.HSMSlotCollection, error) {
	f, ok := filters.(*hsmSlotDBFilter)
	if !ok {
		return nil, errors.Internal().WithMessage("invalid query filters provided")
	}

	if f.Pagination != nil {
		f.Pagination.Limit++
	}
	storageData, err := r.infra.List(ctx, *f.HSMSlotDBFilter)
	if err != nil {
		return nil, mapPersistenceErrorToSignerError(err)
	}

	collection := hsmslot.HSMSlotCollection{}
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
	collection.Items, err = mapSliceFromDB(storageData)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	return &collection, nil
}

// Filter creates a new filter.
func (r *Repository) Filter() hsmslot.HSMSlotFilters {
	storageFilter := hsmSlotDBFilter{
		HSMSlotDBFilter: &hsmslotdb.HSMSlotDBFilter{
			HSMSlotDB: hsmslotdb.HSMSlotDB{},
		},
	}
	return &storageFilter
}

// Repository implementation of hsmslot.HSMSlotStorage
type Repository struct {
	infra *hsmslotdb.HSMSlotRepositoryInfra
}

// RepositoryOptions configures a Repository
type RepositoryOptions struct {
	Infra *hsmslotdb.HSMSlotRepositoryInfra
}

// NewRepository creates a Repository with the given options
func NewRepository(options RepositoryOptions) (*Repository, error) {
	return &Repository{
		infra: options.Infra,
	}, nil
}

var _ hsmslot.HSMSlotFilters = new(hsmSlotDBFilter)

// FilterByApplicationID filters the items by the specified Application ID.
func (filter *hsmSlotDBFilter) FilterByApplicationID(applicationID entities.StandardID) hsmslot.HSMSlotFilters {
	filter.ApplicationID = applicationID.ID
	filter.AppendFilter(postgres.NewEqualFilter("application_id"))
	return filter
}

// FilterByHSMModuleID filters the items by the specified HSM Module ID.
func (filter *hsmSlotDBFilter) FilterByHSMModuleID(hsmID entities.StandardID) hsmslot.HSMSlotFilters {
	filter.HSMModuleID = hsmID.ID
	filter.AppendFilter(postgres.NewEqualFilter("hardware_security_module_id"))
	return filter
}

// OrderByCreationDate orders resources in storage by creation date.
func (filter *hsmSlotDBFilter) OrderByCreationDate(orderDirection persistence.OrderDirection) hsmslot.HSMSlotFilters {
	filter.HSMSlotDBFilter = filter.HSMSlotDBFilter.Sort("creation_date", orderDirection)
	return filter
}

// OrderByLastUpdateDate orders resources in storage by last update date.
func (filter *hsmSlotDBFilter) OrderByLastUpdateDate(orderDirection persistence.OrderDirection) hsmslot.HSMSlotFilters {
	filter.HSMSlotDBFilter = filter.HSMSlotDBFilter.Sort("last_update", orderDirection)
	return filter
}

// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
func (filter *hsmSlotDBFilter) Paged(limit int, offset int) hsmslot.HSMSlotFilters {
	filter.HSMSlotDBFilter = filter.HSMSlotDBFilter.Paged(limit, offset)
	return filter
}

type hsmSlotDBFilter struct {
	*hsmslotdb.HSMSlotDBFilter
}
