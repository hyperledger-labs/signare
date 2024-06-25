package applicationdb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	addApplicationMapperID    = "signare.application.insert"
	allApplicationMapperID    = "signare.application.list"
	getApplicationMapperID    = "signare.application.getById"
	editApplicationMapperID   = "signare.application.update"
	removeApplicationMapperID = "signare.application.delete"
	existsApplicationMapperID = "signare.application.exists"
)

func (repository *ApplicationRepositoryInfra) Add(ctx context.Context, db ApplicationCreateDB) (*ApplicationDB, error) {
	db.ResourceVersion = uuid.NewString()
	err := repository.genericStorage.ExecuteStmt(ctx, addApplicationMapperID, db)
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

func (repository *ApplicationRepositoryInfra) All(ctx context.Context, filters ApplicationDBFilter) ([]ApplicationDB, error) {
	applicationDBItems := make([]ApplicationDB, 0)
	err := repository.genericStorage.QueryAll(ctx, allApplicationMapperID, &filters, &applicationDBItems)
	if err != nil {
		return nil, err
	}
	return applicationDBItems, nil
}

func (repository *ApplicationRepositoryInfra) Get(ctx context.Context, id entities.StandardID) ([]ApplicationDB, error) {
	var applicationDBItems []ApplicationDB
	db := ApplicationDB{}
	db.StandardID = id

	err := repository.genericStorage.QueryAll(ctx, getApplicationMapperID, db, &applicationDBItems)
	if err != nil {
		return nil, err
	}
	return applicationDBItems, nil
}

func (repository *ApplicationRepositoryInfra) Edit(ctx context.Context, db ApplicationUpdateDB) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, existsErr := repository.Exists(ctx, db.StandardID)
	if existsErr != nil {
		return nil, existsErr
	}

	db.NewResourceVersion = uuid.NewString()
	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editApplicationMapperID, db)
}

func (repository *ApplicationRepositoryInfra) Remove(ctx context.Context, id entities.StandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := ApplicationDB{}
	db.StandardID = id

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeApplicationMapperID, db)
}

func (repository *ApplicationRepositoryInfra) Exists(ctx context.Context, id entities.StandardID) ([]ApplicationExistsDB, error) {
	var existsResult []ApplicationExistsDB

	db := ApplicationDB{}
	db.StandardID = id

	err := repository.genericStorage.QueryAll(ctx, existsApplicationMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type ApplicationRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type ApplicationRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideApplicationRepositoryInfra(options ApplicationRepositoryInfraOptions) (*ApplicationRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &ApplicationRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
