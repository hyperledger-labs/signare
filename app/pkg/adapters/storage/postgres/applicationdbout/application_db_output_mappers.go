package applicationdbout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/applicationdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
)

func mapToCreateDB(application application.Application) (*applicationdb.ApplicationCreateDB, error) {
	if len(application.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}
	if len(application.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	db := applicationdb.ApplicationCreateDB{
		ApplicationDB: applicationdb.ApplicationDB{
			StandardID:         application.StandardID,
			InternalResourceID: application.InternalResourceID.String(),
			ChainID:            application.ChainID.String(),
			CreationDate:       application.CreationDate.ToInt64(),
			LastUpdate:         application.LastUpdate.ToInt64(),
		},
	}
	if application.Description != nil {
		db.Description = application.Description
	}
	return &db, nil
}

func mapToUpdateDB(application application.Application) (*applicationdb.ApplicationUpdateDB, error) {
	if len(application.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}

	db := applicationdb.ApplicationUpdateDB{
		ApplicationDB: applicationdb.ApplicationDB{
			StandardID:      application.StandardID,
			ChainID:         application.ChainID.String(),
			LastUpdate:      application.LastUpdate.ToInt64(),
			ResourceVersion: application.ResourceVersion,
		},
	}
	if application.Description != nil {
		db.Description = application.Description
	}
	return &db, nil
}

func mapFromDB(db applicationdb.ApplicationDB) (*application.Application, error) {
	if len(db.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	chainID, err := entities.NewInt256FromString(db.ChainID)
	if err != nil {
		return nil, errors.Internal().WithMessage("chainID [%s] can not be casted to Int256", db.ChainID)
	}

	description := ""
	if db.Description != nil {
		description = *db.Description
	}

	app := application.Application{
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
		ChainID:            *chainID,
		Description:        &description,
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
	}
	return &app, nil
}

func mapSliceFromDB(dbSlice []applicationdb.ApplicationDB) ([]application.Application, error) {
	applicationSlice := make([]application.Application, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		applicationSlice[index] = *item
	}

	return applicationSlice, nil
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
