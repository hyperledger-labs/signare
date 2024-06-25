package hsmmoduledb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	addHardwareSecurityModuleMapperID    = "signare.hardwareSecurityModule.insert"
	getHardwareSecurityModuleMapperID    = "signare.hardwareSecurityModule.getById"
	removeHardwareSecurityModuleMapperID = "signare.hardwareSecurityModule.delete"
	editHardwareSecurityModuleMapperID   = "signare.hardwareSecurityModule.update"
	listHardwareSecurityModulesMapperID  = "signare.hardwareSecurityModule.list"
	existsHardwareSecurityModuleMapperID = "signare.hardwareSecurityModule.exists"
)

func (repository *HardwareSecurityModuleRepositoryInfra) Add(ctx context.Context, db HardwareSecurityModuleCreateDB) (*HardwareSecurityModuleDB, error) {
	db.ResourceVersion = uuid.NewString()
	err := repository.genericStorage.ExecuteStmt(ctx, addHardwareSecurityModuleMapperID, db)
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

func (repository *HardwareSecurityModuleRepositoryInfra) Get(ctx context.Context, input entities.StandardID) ([]HardwareSecurityModuleDB, error) {
	var hardwareSecurityModuleDBItems []HardwareSecurityModuleDB
	db := HardwareSecurityModuleDB{}
	db.ID = input.ID

	err := repository.genericStorage.QueryAll(ctx, getHardwareSecurityModuleMapperID, db, &hardwareSecurityModuleDBItems)
	if err != nil {
		return nil, err
	}
	return hardwareSecurityModuleDBItems, nil
}

func (repository *HardwareSecurityModuleRepositoryInfra) Edit(ctx context.Context, db HardwareSecurityModuleUpdateDB) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, err := repository.Exists(ctx, db.StandardID)
	if err != nil {
		return nil, err
	}

	db.NewResourceVersion = uuid.NewString()

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editHardwareSecurityModuleMapperID, db)
}

func (repository *HardwareSecurityModuleRepositoryInfra) Remove(ctx context.Context, input entities.StandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := HardwareSecurityModuleDB{}
	db.ID = input.ID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeHardwareSecurityModuleMapperID, db)
}

func (repository *HardwareSecurityModuleRepositoryInfra) All(ctx context.Context, filters HardwareSecurityModuleDBFilter) ([]HardwareSecurityModuleDB, error) {
	hardwareSecurityModuleDBItems := make([]HardwareSecurityModuleDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listHardwareSecurityModulesMapperID, &filters, &hardwareSecurityModuleDBItems)
	if err != nil {
		return nil, err
	}
	return hardwareSecurityModuleDBItems, nil
}

func (repository *HardwareSecurityModuleRepositoryInfra) Exists(ctx context.Context, id entities.StandardID) ([]HardwareSecurityModuleExistsDB, error) {
	var existsResult []HardwareSecurityModuleExistsDB

	db := HardwareSecurityModuleDB{}
	db.StandardID = id

	err := repository.genericStorage.QueryAll(ctx, existsHardwareSecurityModuleMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type HardwareSecurityModuleRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type HardwareSecurityModuleRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideHardwareSecurityModuleRepositoryInfra(options HardwareSecurityModuleRepositoryInfraOptions) (*HardwareSecurityModuleRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &HardwareSecurityModuleRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
