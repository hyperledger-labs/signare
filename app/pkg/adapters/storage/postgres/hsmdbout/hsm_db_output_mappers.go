package hsmdbout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmmoduledb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
)

func mapToCreateDB(module hsmmodule.HSMModule) (*hsmmoduledb.HardwareSecurityModuleCreateDB, error) {
	if len(module.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}
	if len(module.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	dbModuleKind, mapErr := mapDBModuleKindFrom(module.Kind)
	if mapErr != nil {
		return nil, mapErr
	}
	configuration := mapConfigurationDBFrom(module)
	db := hsmmoduledb.HardwareSecurityModuleCreateDB{
		HardwareSecurityModuleDB: hsmmoduledb.HardwareSecurityModuleDB{
			StandardID:         module.StandardID,
			InternalResourceID: module.InternalResourceID.String(),
			Kind:               *dbModuleKind,
			Configuration:      *configuration,
			LastUpdate:         module.LastUpdate.ToInt64(),
			CreationDate:       module.CreationDate.ToInt64(),
		},
	}
	if module.Description == nil {
		db.Description = ""
	} else {
		db.Description = *module.Description
	}

	return &db, nil
}

func mapToUpdateDB(module hsmmodule.HSMModule) (*hsmmoduledb.HardwareSecurityModuleUpdateDB, error) {
	if len(module.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}
	dbModuleKind, mapErr := mapDBModuleKindFrom(module.Kind)
	if mapErr != nil {
		return nil, mapErr
	}
	configuration := mapConfigurationDBFrom(module)

	db := hsmmoduledb.HardwareSecurityModuleUpdateDB{
		HardwareSecurityModuleDB: hsmmoduledb.HardwareSecurityModuleDB{
			StandardID:         module.StandardID,
			Kind:               *dbModuleKind,
			Configuration:      *configuration,
			InternalResourceID: module.InternalResourceID.String(),
			CreationDate:       module.CreationDate.ToInt64(),
			LastUpdate:         module.LastUpdate.ToInt64(),
			ResourceVersion:    module.ResourceVersion,
		},
	}
	if module.Description == nil {
		db.Description = ""
	} else {
		db.Description = *module.Description
	}
	return &db, nil
}

func mapFromDB(db hsmmoduledb.HardwareSecurityModuleDB) (*hsmmodule.HSMModule, error) {
	if len(db.ID) == 0 {
		return nil, errors.Internal().WithMessage("id cannot be empty")
	}
	if len(db.InternalResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'InternalResourceID' cannot be empty")
	}
	configuration := mapUseCaseConfiguration(db)
	useCaseHSMType, mapErr := mapUseCaseModuleKindFrom(db.Kind)
	if mapErr != nil {
		return nil, mapErr
	}

	hsmModule := hsmmodule.HSMModule{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: db.StandardID,
				Timestamps: entities.Timestamps{
					CreationDate: time.TimestampFromInt64(db.CreationDate),
					LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
				},
			},
			ResourceVersion: db.ResourceVersion,
		},
		InternalResourceID: entities.InternalResourceID(db.InternalResourceID),
		Description:        &db.Description,
		Kind:               *useCaseHSMType,
		Configuration:      *configuration,
	}
	return &hsmModule, nil
}

func mapSliceFromDB(dbSlice []hsmmoduledb.HardwareSecurityModuleDB) ([]hsmmodule.HSMModule, error) {
	hsmModuleSlice := make([]hsmmodule.HSMModule, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		hsmModuleSlice[index] = *item
	}

	return hsmModuleSlice, nil
}

func mapUseCaseConfiguration(module hsmmoduledb.HardwareSecurityModuleDB) *hsmmodule.HSMModuleConfiguration {
	var configuration hsmmodule.HSMModuleConfiguration
	if module.Kind == string(hsmmodule.SoftHSMModuleKind) {
		configuration.SoftHSMConfiguration = &hsmmodule.SoftHSMConfiguration{}
	}

	return &configuration
}

func mapConfigurationDBFrom(module hsmmodule.HSMModule) *string {
	var configuration string
	if module.Kind == hsmmodule.SoftHSMModuleKind {
		// SoftHSM configuration is static and not persisted
		configuration = ""
	}

	return &configuration
}

func mapUseCaseModuleKindFrom(kind string) (*hsmmodule.ModuleKind, error) {
	if kind == hsmmoduledb.SoftHSMModuleKind {
		k := hsmmodule.SoftHSMModuleKind
		return &k, nil
	}
	return nil, errors.Internal().WithMessage("couldn't map '%s' to usecase HSM kind", kind)
}

func mapDBModuleKindFrom(moduleKind hsmmodule.ModuleKind) (*string, error) {
	if moduleKind == hsmmodule.SoftHSMModuleKind {
		kind := hsmmoduledb.SoftHSMModuleKind
		return &kind, nil
	}
	return nil, errors.InvalidArgument().WithMessage("couldn't map '%s' to database HSM type", moduleKind)
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
