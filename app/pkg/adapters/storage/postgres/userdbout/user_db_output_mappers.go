package userdbout

import (
	"encoding/json"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/userdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

func mapToCreateDB(user user.User) (*userdb.UserCreateDB, error) {
	if len(user.ID) == 0 {
		return nil, errors.Internal().WithMessage("'AdminID' cannot be empty")
	}
	if len(user.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if len(user.ApplicationID) == 0 {
		return nil, errors.Internal().WithMessage("'ApplicationID' cannot be empty")
	}
	if len(user.Roles) == 0 {
		return nil, errors.Internal().WithMessage("'Roles' cannot be empty")
	}
	roles, err := json.Marshal(user.Roles)
	if err != nil {
		return nil, err
	}

	db := userdb.UserCreateDB{
		UserDB: userdb.UserDB{
			ApplicationStandardID: user.ApplicationStandardID,
			InternalResourceID:    user.InternalResourceID.String(),
			Roles:                 string(roles),
			CreationDate:          user.CreationDate.ToInt64(),
			LastUpdate:            user.LastUpdate.ToInt64(),
		},
	}
	if user.Description != nil {
		db.Description = *user.Description
	}
	return &db, nil
}

func mapToUpdateDB(user user.User) (*userdb.UserUpdateDB, error) {
	if len(user.ID) == 0 {
		return nil, errors.Internal().WithMessage("'AdminID' cannot be empty")
	}
	if len(user.ApplicationID) == 0 {
		return nil, errors.Internal().WithMessage("'ApplicationID' cannot be empty")
	}
	if len(user.Roles) == 0 {
		return nil, errors.Internal().WithMessage("'Roles' cannot be empty")
	}
	roles, err := json.Marshal(user.Roles)
	if err != nil {
		return nil, err
	}

	db := userdb.UserUpdateDB{
		UserDB: userdb.UserDB{
			ApplicationStandardID: user.ApplicationStandardID,
			Roles:                 string(roles),
			CreationDate:          user.CreationDate.ToInt64(),
			InternalResourceID:    user.InternalResourceID.String(),
			LastUpdate:            user.LastUpdate.ToInt64(),
			ResourceVersion:       user.ResourceVersion,
		},
	}
	if user.Description != nil {
		db.Description = *user.Description
	}
	return &db, nil
}

func mapFromDB(db userdb.UserDB) (*user.User, error) {
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	var roles []string

	err := json.Unmarshal([]byte(db.Roles), &roles)
	if err != nil {
		return nil, err
	}

	return &user.User{
		ApplicationStandardResourceMeta: entities.ApplicationStandardResourceMeta{
			ApplicationStandardResource: entities.ApplicationStandardResource{
				ApplicationStandardID: entities.ApplicationStandardID{
					ID:            db.ID,
					ApplicationID: db.ApplicationID,
				},
				Timestamps: entities.Timestamps{
					CreationDate: time.TimestampFromInt64(db.CreationDate),
					LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
				},
			},
			ResourceVersion: db.ResourceVersion,
		},
		Roles:              roles,
		Description:        &db.Description,
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
	}, nil
}

func mapSliceFromDB(dbSlice []userdb.UserDB) ([]user.User, error) {
	userSlice := make([]user.User, len(dbSlice))
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
