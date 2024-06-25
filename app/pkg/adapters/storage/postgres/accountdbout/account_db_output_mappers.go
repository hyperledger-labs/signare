package accountdbout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/accountdb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

func mapToCreateDB(account user.Account) (*accountdb.AccountCreateDB, error) {
	if len(account.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if account.Address.IsEmpty() {
		return nil, errors.Internal().WithMessage("'Address' cannot be empty")
	}
	if len(account.UserID) == 0 {
		return nil, errors.Internal().WithMessage("'UserID' cannot be empty")
	}
	if len(account.ApplicationID) == 0 {
		return nil, errors.Internal().WithMessage("'ApplicationID' cannot be empty")
	}

	db := accountdb.AccountCreateDB{
		AccountDB: accountdb.AccountDB{
			InternalResourceID: account.InternalResourceID.String(),
			Address:            account.Address.String(),
			ApplicationID:      account.ApplicationID,
			UserID:             account.UserID,
			CreationDate:       account.CreationDate.ToInt64(),
			LastUpdate:         account.LastUpdate.ToInt64(),
		},
	}
	return &db, nil
}

func mapFromDB(db accountdb.AccountDB) (*user.Account, error) {
	addr, err := address.NewFromHexString(db.Address)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	if addr.IsEmpty() {
		return nil, errors.Internal().WithMessage("'Address' cannot be empty")
	}
	if len(db.UserID) == 0 {
		return nil, errors.Internal().WithMessage("'UserID' cannot be empty")
	}
	if len(db.ApplicationID) == 0 {
		return nil, errors.Internal().WithMessage("'ApplicationID' cannot be empty")
	}

	return &user.Account{
		AccountID: user.AccountID{
			Address:       addr,
			UserID:        db.UserID,
			ApplicationID: db.ApplicationID,
		},
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
		Timestamps: entities.Timestamps{
			CreationDate: time.TimestampFromInt64(db.CreationDate),
			LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
		},
	}, nil
}

func mapToAccountID(id user.AccountID) accountdb.AccountID {
	return accountdb.AccountID{
		Address:       id.Address.String(),
		ApplicationID: id.ApplicationID,
		UserID:        id.UserID,
	}
}

func mapToAccountApplicationAddresses(applicationID string, address address.Address) accountdb.AccountApplicationAddresses {
	return accountdb.AccountApplicationAddresses{
		Address:       address.String(),
		ApplicationID: applicationID,
	}
}

func mapSliceFromDB(dbSlice []accountdb.AccountDB) ([]user.Account, error) {
	accountSlice := make([]user.Account, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		accountSlice[index] = *item
	}

	return accountSlice, nil
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
