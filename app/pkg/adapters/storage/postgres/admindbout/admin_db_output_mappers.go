package admindbout

import (
	"encoding/json"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/admindb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
)

func mapToCreateDB(admin admin.Admin) (*admindb.AdminCreateDB, error) {
	if len(admin.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(admin.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if len(admin.Roles) == 0 {
		return nil, errors.Internal().WithMessage("'Roles' cannot be empty")
	}
	roles, err := json.Marshal(admin.Roles)
	if err != nil {
		return nil, err
	}
	var description string
	if admin.Description != nil {
		description = *admin.Description
	}

	db := admindb.AdminCreateDB{
		AdminDB: admindb.AdminDB{
			StandardID:         admin.StandardID,
			InternalResourceID: admin.InternalResourceID.String(),
			Roles:              string(roles),
			Description:        description,
			CreationDate:       admin.CreationDate.ToInt64(),
			LastUpdate:         admin.LastUpdate.ToInt64(),
		},
	}
	if admin.Description != nil {
		db.Description = *admin.Description
	}
	return &db, nil
}

func mapToUpdateDB(admin admin.Admin) (*admindb.AdminUpdateDB, error) {
	if len(admin.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(admin.Roles) == 0 {
		return nil, errors.Internal().WithMessage("'Roles' cannot be empty")
	}
	roles, err := json.Marshal(admin.Roles)
	if err != nil {
		return nil, err
	}
	var description string
	if admin.Description != nil {
		description = *admin.Description
	}

	db := admindb.AdminUpdateDB{
		AdminDB: admindb.AdminDB{
			StandardID:      admin.StandardID,
			Roles:           string(roles),
			Description:     description,
			CreationDate:    admin.CreationDate.ToInt64(),
			LastUpdate:      admin.LastUpdate.ToInt64(),
			ResourceVersion: admin.ResourceVersion,
		},
	}
	if admin.Description != nil {
		db.Description = *admin.Description
	}
	return &db, nil
}

func mapFromDB(db admindb.AdminDB) (*admin.Admin, error) {
	if len(db.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if len(db.Roles) == 0 {
		return nil, errors.Internal().WithMessage("'Roles' cannot be empty")
	}

	var roles []string
	err := json.Unmarshal([]byte(db.Roles), &roles)
	if err != nil {
		return nil, err
	}

	return &admin.Admin{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: entities.StandardID{
					ID: db.ID,
				},
				Timestamps: entities.Timestamps{
					CreationDate: time.TimestampFromInt64(db.CreationDate),
					LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
				},
			},
			ResourceVersion: db.ResourceVersion,
		},
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
		Description:        &db.Description,
		Roles:              roles,
	}, nil
}

func mapSliceFromDB(dbSlice []admindb.AdminDB) ([]admin.Admin, error) {
	adminSlice := make([]admin.Admin, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		adminSlice[index] = *item
	}

	return adminSlice, nil
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
