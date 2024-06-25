package accountdb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
)

const (
	addAccountMapperID                  = "signare.account.insert"
	getAccountMapperID                  = "signare.account.getById"
	removeAccountMapperID               = "signare.account.delete"
	listAccountsMapperID                = "signare.account.list"
	removeAllAccountsForAddressMapperID = "signare.account.deleteAllForAddress"
)

func (repository *AccountRepositoryInfra) Add(ctx context.Context, db AccountCreateDB) (*AccountDB, error) {
	err := repository.genericStorage.ExecuteStmt(ctx, addAccountMapperID, db)
	if err != nil {
		return nil, err
	}

	result, err := repository.Get(ctx, AccountID{
		Address:       db.Address,
		ApplicationID: db.ApplicationID,
		UserID:        db.UserID,
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, persistence.NewEntryNotAddedError()
	}

	return &result[0], nil
}

func (repository *AccountRepositoryInfra) Get(ctx context.Context, input AccountID) ([]AccountDB, error) {
	var accountDBItems []AccountDB
	db := AccountDB{}
	db.Address = input.Address
	db.ApplicationID = input.ApplicationID
	db.UserID = input.UserID

	err := repository.genericStorage.QueryAll(ctx, getAccountMapperID, db, &accountDBItems)
	if err != nil {
		return nil, err
	}
	return accountDBItems, nil
}

func (repository *AccountRepositoryInfra) Remove(ctx context.Context, input AccountID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := AccountDB{}
	db.Address = input.Address
	db.ApplicationID = input.ApplicationID
	db.UserID = input.UserID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeAccountMapperID, db)
}

func (repository *AccountRepositoryInfra) List(ctx context.Context, filters AccountDBFilter) ([]AccountDB, error) {
	accountDBItems := make([]AccountDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listAccountsMapperID, &filters, &accountDBItems)
	if err != nil {
		return nil, err
	}
	return accountDBItems, nil
}

func (repository *AccountRepositoryInfra) RemoveAllForAddress(ctx context.Context, input AccountApplicationAddresses) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := AccountDB{}
	db.Address = input.Address
	db.ApplicationID = input.ApplicationID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeAllAccountsForAddressMapperID, db)
}

type AccountRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type AccountRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideAccountRepositoryInfra(options AccountRepositoryInfraOptions) (*AccountRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &AccountRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
