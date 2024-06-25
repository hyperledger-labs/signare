package admindb

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	addAdminMapperID    = "signare.admin.insert"
	getAdminMapperID    = "signare.admin.getById"
	editAdminMapperID   = "signare.admin.update"
	removeAdminMapperID = "signare.admin.delete"
	listAdminsMapperID  = "signare.admin.list"
	existsAdminMapperID = "signare.admin.exists"
)

func (repository *AdminRepositoryInfra) Add(ctx context.Context, db AdminCreateDB) (*AdminDB, error) {
	db.ResourceVersion = uuid.NewString()
	err := repository.genericStorage.ExecuteStmt(ctx, addAdminMapperID, db)
	if err != nil {
		return nil, err
	}

	result, err := repository.Get(ctx, db.StandardID)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, persistence.NewEntryNotAddedError()
	}

	return &result[0], nil
}

func (repository *AdminRepositoryInfra) Get(ctx context.Context, input entities.StandardID) ([]AdminDB, error) {
	var adminDBItems []AdminDB
	db := AdminDB{}
	db.ID = input.ID

	err := repository.genericStorage.QueryAll(ctx, getAdminMapperID, db, &adminDBItems)
	if err != nil {
		return nil, err
	}
	return adminDBItems, nil
}

func (repository *AdminRepositoryInfra) Edit(ctx context.Context, db AdminUpdateDB) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, existsErr := repository.Exists(ctx, db.StandardID)
	if existsErr != nil {
		return nil, existsErr
	}

	db.NewResourceVersion = uuid.NewString()
	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editAdminMapperID, db)
}

func (repository *AdminRepositoryInfra) Remove(ctx context.Context, input entities.StandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := AdminDB{}
	db.ID = input.ID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeAdminMapperID, db)
}

func (repository *AdminRepositoryInfra) List(ctx context.Context, filters AdminDBFilter) ([]AdminDB, error) {
	adminDBItems := make([]AdminDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listAdminsMapperID, &filters, &adminDBItems)
	if err != nil {
		return nil, err
	}
	return adminDBItems, nil
}

func (repository *AdminRepositoryInfra) Exists(ctx context.Context, id entities.StandardID) ([]AdminExistsDB, error) {
	var existsResult []AdminExistsDB

	db := AdminDB{}
	db.StandardID = id

	err := repository.genericStorage.QueryAll(ctx, existsAdminMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type AdminRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type AdminRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideAdminRepositoryInfra(options AdminRepositoryInfraOptions) *AdminRepositoryInfra {
	return &AdminRepositoryInfra{
		genericStorage: options.GenericStorage,
	}
}
