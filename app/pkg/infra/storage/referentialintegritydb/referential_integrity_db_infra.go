package referentialintegritydb

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

const (
	addReferentialIntegrityEntryMapperID                     = "signare.referentialIntegrityEntry.insert"
	getReferentialIntegrityEntryMapperID                     = "signare.referentialIntegrityEntry.getById"
	removeReferentialIntegrityEntryMapperID                  = "signare.referentialIntegrityEntry.delete"
	listReferentialIntegrityEntriesMapperID                  = "signare.referentialIntegrityEntry.list"
	removeAllFromResourceReferentialIntegrityEntriesMapperID = "signare.referentialIntegrityEntry.deleteAllFromResource"
)

func (repository *ReferentialIntegrityEntryRepositoryInfra) Add(ctx context.Context, db ReferentialIntegrityEntryCreateDB) (*ReferentialIntegrityEntryDB, error) {
	err := repository.genericStorage.ExecuteStmt(ctx, addReferentialIntegrityEntryMapperID, db)
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

func (repository *ReferentialIntegrityEntryRepositoryInfra) Get(ctx context.Context, input entities.StandardID) ([]ReferentialIntegrityEntryDB, error) {
	var hardwareSecurityModuleDBItems []ReferentialIntegrityEntryDB
	db := ReferentialIntegrityEntryDB{}
	db.ID = input.ID

	err := repository.genericStorage.QueryAll(ctx, getReferentialIntegrityEntryMapperID, db, &hardwareSecurityModuleDBItems)
	if err != nil {
		return nil, err
	}
	return hardwareSecurityModuleDBItems, nil
}

func (repository *ReferentialIntegrityEntryRepositoryInfra) Remove(ctx context.Context, input entities.StandardID) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := ReferentialIntegrityEntryDB{}
	db.ID = input.ID

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeReferentialIntegrityEntryMapperID, db)
}

func (repository *ReferentialIntegrityEntryRepositoryInfra) RemoveAllFromResource(ctx context.Context, resourceID, resourceKind string) (*persistence.ExecuteStmtWithStorageResultOutput, error) {
	db := ReferentialIntegrityEntryDB{}
	db.ResourceID = resourceID
	db.ResourceKind = resourceKind

	return repository.genericStorage.ExecuteStmtWithStorageResult(ctx, removeAllFromResourceReferentialIntegrityEntriesMapperID, db)
}

func (repository *ReferentialIntegrityEntryRepositoryInfra) All(ctx context.Context, filters ReferentialIntegrityEntryDBFilter) ([]ReferentialIntegrityEntryDB, error) {
	referentialIntegrityEntryDBItems := make([]ReferentialIntegrityEntryDB, 0)
	err := repository.genericStorage.QueryAll(ctx, listReferentialIntegrityEntriesMapperID, &filters, &referentialIntegrityEntryDBItems)
	if err != nil {
		return nil, err
	}
	return referentialIntegrityEntryDBItems, nil
}

type ReferentialIntegrityEntryRepositoryInfraOptions struct {
	GenericStorage persistence.Storage
}

type ReferentialIntegrityEntryRepositoryInfra struct {
	genericStorage persistence.Storage
}

func ProvideReferentialIntegrityEntryRepositoryInfra(options ReferentialIntegrityEntryRepositoryInfraOptions) (*ReferentialIntegrityEntryRepositoryInfra, error) {
	if options.GenericStorage == nil {
		return nil, fmt.Errorf("mandatory 'GenericStorage' not provided")
	}
	return &ReferentialIntegrityEntryRepositoryInfra{
		genericStorage: options.GenericStorage,
	}, nil
}
