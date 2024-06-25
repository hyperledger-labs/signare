package hsmslotdbout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmslotdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
)

func mapToCreateDB(slot hsmslot.HSMSlot) (*hsmslotdb.HSMSlotCreateDB, error) {
	if len(slot.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(slot.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if len(slot.ApplicationID) == 0 {
		return nil, errors.Internal().WithMessage("'ApplicationID' cannot be empty")
	}
	if len(slot.HSMModuleID) == 0 {
		return nil, errors.Internal().WithMessage("'HSMModuleID' cannot be empty")
	}
	if len(slot.Slot) == 0 {
		return nil, errors.Internal().WithMessage("'Slot' cannot be empty")
	}
	if len(slot.Pin) == 0 {
		return nil, errors.Internal().WithMessage("'Pin' cannot be empty")
	}
	return &hsmslotdb.HSMSlotCreateDB{
		HSMSlotDB: hsmslotdb.HSMSlotDB{
			StandardID:         slot.StandardID,
			InternalResourceID: slot.InternalResourceID.String(),
			ApplicationID:      slot.ApplicationID,
			HSMModuleID:        slot.HSMModuleID,
			Slot:               slot.Slot,
			Pin:                slot.Pin,
			CreationDate:       slot.CreationDate.ToInt64(),
			LastUpdate:         slot.LastUpdate.ToInt64(),
		},
	}, nil
}

func mapToUpdatePinDB(slot hsmslot.HSMSlot) (*hsmslotdb.HSMSlotUpdatePinDB, error) {
	if len(slot.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(slot.Pin) == 0 {
		return nil, errors.Internal().WithMessage("'Pin' cannot be empty")
	}
	return &hsmslotdb.HSMSlotUpdatePinDB{
		StandardID:      slot.StandardID,
		ResourceVersion: slot.ResourceVersion,
		Pin:             slot.Pin,
		LastUpdate:      slot.LastUpdate.ToInt64(),
	}, nil
}

func mapFromDB(db hsmslotdb.HSMSlotDB) (*hsmslot.HSMSlot, error) {
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	return &hsmslot.HSMSlot{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: db.StandardID,
				Timestamps: entities.Timestamps{
					CreationDate: time.TimestampFromInt64(db.CreationDate),
					LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
				},
			},
			ResourceVersion: db.ResourceVersion,
		},
		ApplicationID:      db.ApplicationID,
		HSMModuleID:        db.HSMModuleID,
		Slot:               db.Slot,
		Pin:                db.Pin,
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
	}, nil
}

func mapSliceFromDB(dbSlice []hsmslotdb.HSMSlotDB) ([]hsmslot.HSMSlot, error) {
	userSlice := make([]hsmslot.HSMSlot, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		userSlice[index] = *item
	}
	return userSlice, nil
}

func mapPersistenceErrorToSignerError(err error) error {
	if persistence.IsAlreadyExists(err) {
		return errors.AlreadyExistsFromErr(err)
	}
	if persistence.IsNotFound(err) {
		return errors.NotFoundFromErr(err)
	}
	if persistence.IsEntryNotAdded(err) {
		return errors.InternalFromErr(err)
	}
	return errors.InternalFromErr(err)
}
