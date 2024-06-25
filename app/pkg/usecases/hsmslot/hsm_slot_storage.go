package hsmslot

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// HSMSlotStorage defines the functionality to interact with the HSMSlot in storage.
type HSMSlotStorage interface {
	// Add an HSMSlot in storage.
	Add(ctx context.Context, data HSMSlot) (*HSMSlot, error)
	// Get an HSMSlot from storage.
	Get(ctx context.Context, id entities.StandardID) (*HSMSlot, error)
	// GetByApplication an HSMSlot from storage.
	GetByApplication(ctx context.Context, applicationID entities.StandardID) (*HSMSlot, error)
	// EditPin of an HSMSlot in storage.
	EditPin(ctx context.Context, data HSMSlot) (*HSMSlot, error)
	// Remove an HSMSlot in storage.
	Remove(ctx context.Context, id entities.StandardID) (*HSMSlot, error)
	// All HSMSlot in storage.
	All(ctx context.Context, filters HSMSlotFilters) (*HSMSlotCollection, error)
	// Filter all HSMSlot by specific properties.
	Filter() HSMSlotFilters
}

// HSMSlotFilters defines filter options for retrieving HSMSlots from storage.
type HSMSlotFilters interface {
	// FilterByApplicationID filters the slots for a specific Application.
	FilterByApplicationID(applicationID entities.StandardID) HSMSlotFilters
	// FilterByHSMModuleID filters the slots for a specific HSM.
	FilterByHSMModuleID(hsmID entities.StandardID) HSMSlotFilters
	// OrderByCreationDate orders HSMSlot in storage by creation date.
	OrderByCreationDate(orderDirection persistence.OrderDirection) HSMSlotFilters
	// OrderByLastUpdateDate orders HSMSlot in storage by last update date.
	OrderByLastUpdateDate(orderDirection persistence.OrderDirection) HSMSlotFilters
	// Paged limits the maximum amount of items to limit parameter and starts the list in offset parameter.
	Paged(limit int, offset int) HSMSlotFilters
}
