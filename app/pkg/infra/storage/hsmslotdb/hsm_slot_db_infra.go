package hsmslotdb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/google/uuid"
)

const (
	addSlotMapperID              = "signare.hardwareSecurityModuleSlot.insert"
	getSlotMapperID              = "signare.hardwareSecurityModuleSlot.getById"
	getSlotByApplicationMapperID = "signare.hardwareSecurityModuleSlot.getByApplication"
	editPinSlotMapperID          = "signare.hardwareSecurityModuleSlot.updatePin"
	removeSlotMapperID           = "signare.hardwareSecurityModuleSlot.delete"
	listSlotMapperID             = "signare.hardwareSecurityModuleSlot.list"
	exitsSlotMapperID            = "signare.hardwareSecurityModuleSlot.exists"
)

func (repository *HSMSlotRepositoryInfra) Add(ctx context.Context, db HSMSlotCreateDB) (*HSMSlotDB, error) {
	db.ResourceVersion = uuid.NewString()
	err := repository.genericStorage.ExecuteStmt(ctx, addSlotMapperID, db)
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

func (repository *HSMSlotRepositoryInfra) Get(ctx context.Context, id entities.StandardID) ([]HSMSlotDB, error) {
	var slotDBItems []HSMSlotDB
	db := HSMSlotDB{}
	db.ID = id.ID

	err := repository.genericStorage.QueryAll(ctx, getSlotMapperID, db, &slotDBItems)
	if err != nil {
		return nil, err
	}
	return slotDBItems, nil
}

func (repository *HSMSlotRepositoryInfra) GetByApplicationID(ctx context.Context, applicationID entities.StandardID) ([]HSMSlotDB, error) {
	var slotDBItems []HSMSlotDB
	db := HSMSlotDB{}
	db.ApplicationID = applicationID.ID

	err := repository.genericStorage.QueryAll(ctx, getSlotByApplicationMapperID, db, &slotDBItems)
	if err != nil {
		return nil, err
	}
	return slotDBItems, nil
}

func (repository *HSMSlotRepositoryInfra) EditPin(ctx context.Context, db HSMSlotUpdatePinDB) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	_, err := repository.Exists(ctx, db.StandardID)
	if err != nil {
		return nil, err
	}

	db.NewResourceVersion = uuid.NewString()
	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, editPinSlotMapperID, db)
}

func (repository *HSMSlotRepositoryInfra) Remove(ctx context.Context, id entities.StandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := HSMSlotDB{}
	db.ID = id.ID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeSlotMapperID, db)
}

func (repository *HSMSlotRepositoryInfra) List(ctx context.Context, filters HSMSlotDBFilter) ([]HSMSlotDB, error) {
	slotDBItems := make([]HSMSlotDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listSlotMapperID, &filters, &slotDBItems)
	if err != nil {
		return nil, err
	}
	return slotDBItems, nil
}

func (repository *HSMSlotRepositoryInfra) Exists(ctx context.Context, id entities.StandardID) ([]HSMSlotExistsDB, error) {
	var existsResult []HSMSlotExistsDB

	db := HSMSlotDB{}
	db.ID = id.ID

	err := repository.genericStorage.QueryAll(ctx, exitsSlotMapperID, db, &existsResult)
	if err != nil {
		return existsResult, err
	}

	if len(existsResult) == 0 || !existsResult[0].Exists {
		return nil, persistence.NewNotFoundError()
	}

	return existsResult, nil
}

type HSMSlotRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type HSMSlotRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideHSMSlotRepositoryInfra(options HSMSlotRepositoryInfraOptions) (*HSMSlotRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &HSMSlotRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
