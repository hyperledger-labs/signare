package userdb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	addUserMapperID    = "signare.user.insert"
	getUserMapperID    = "signare.user.getById"
	editUserMapperID   = "signare.user.update"
	removeUserMapperID = "signare.user.delete"
	listUsersMapperID  = "signare.user.list"
	existsUserMapperID = "signare.user.exists"
)

func (repository *UserRepositoryInfra) Add(ctx context.Context, db UserCreateDB) (*UserDB, error) {
	db.ResourceVersion = uuid.NewString()
	err := repository.genericStorage.ExecuteStmt(ctx, addUserMapperID, db)
	if err != nil {
		return nil, err
	}

	result, err := repository.Get(ctx, db.ApplicationStandardID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, persistence.NewEntryNotAddedError()
	}

	return &result[0], nil
}

func (repository *UserRepositoryInfra) Get(ctx context.Context, id entities.ApplicationStandardID) ([]UserDB, error) {
	var userDBItems []UserDB
	db := UserDB{}
	db.ID = id.ID
	db.ApplicationID = id.ApplicationID

	err := repository.genericStorage.QueryAll(ctx, getUserMapperID, db, &userDBItems)
	if err != nil {
		return nil, err
	}
	return userDBItems, nil
}

func (repository *UserRepositoryInfra) Edit(ctx context.Context, db UserUpdateDB) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, err := repository.Exists(ctx, db.ApplicationStandardID)
	if err != nil {
		return nil, err
	}

	db.NewResourceVersion = uuid.NewString()

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editUserMapperID, db)
}

func (repository *UserRepositoryInfra) Remove(ctx context.Context, id entities.ApplicationStandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := UserDB{}
	db.ID = id.ID
	db.ApplicationID = id.ApplicationID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeUserMapperID, db)
}

func (repository *UserRepositoryInfra) List(ctx context.Context, filters UserDBFilter) ([]UserDB, error) {
	userDBItems := make([]UserDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listUsersMapperID, &filters, &userDBItems)
	if err != nil {
		return nil, err
	}
	return userDBItems, nil
}

func (repository *UserRepositoryInfra) Exists(ctx context.Context, id entities.ApplicationStandardID) ([]UserExistsDB, error) {
	var existsResult []UserExistsDB

	db := UserDB{}
	db.ApplicationStandardID = id

	err := repository.genericStorage.QueryAll(ctx, existsUserMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type UserRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type UserRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideUserRepositoryInfra(options UserRepositoryInfraOptions) (*UserRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &UserRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
