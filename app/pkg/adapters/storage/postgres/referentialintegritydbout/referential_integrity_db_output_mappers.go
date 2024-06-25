package referentialintegritydbout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/referentialintegritydb"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
)

func mapToCreateDB(entry referentialintegrity.ReferentialIntegrityEntry) (*referentialintegritydb.ReferentialIntegrityEntryCreateDB, error) {
	if len(entry.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(entry.ResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'ResourceID' cannot be empty")
	}
	if len(entry.ResourceKind) == 0 {
		return nil, errors.Internal().WithMessage("'ResourceKind' cannot be empty")
	}
	if len(entry.ParentResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'ParentResourceID' cannot be empty")
	}
	if len(entry.ParentResourceKind) == 0 {
		return nil, errors.Internal().WithMessage("'ParentResourceKind' cannot be empty")
	}
	ResourceKind, mapErr := mapTypeToDB(entry.ResourceKind)
	if mapErr != nil {
		return nil, mapErr
	}
	parentResourceKind, mapErr := mapTypeToDB(entry.ParentResourceKind)
	if mapErr != nil {
		return nil, mapErr
	}

	return &referentialintegritydb.ReferentialIntegrityEntryCreateDB{
		ReferentialIntegrityEntryDB: referentialintegritydb.ReferentialIntegrityEntryDB{
			StandardID:         entry.StandardID,
			ResourceID:         entry.ResourceID,
			ResourceKind:       *ResourceKind,
			ParentResourceID:   entry.ParentResourceID,
			ParentResourceKind: *parentResourceKind,
			CreationDate:       entry.CreationDate.ToInt64(),
			LastUpdate:         entry.LastUpdate.ToInt64(),
		},
	}, nil
}

func mapFromDB(db referentialintegritydb.ReferentialIntegrityEntryDB) (*referentialintegrity.ReferentialIntegrityEntry, error) {
	if len(db.ID) == 0 {
		return nil, errors.Internal().WithMessage("'ID' cannot be empty")
	}
	if len(db.ResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'ResourceID' cannot be empty")
	}
	if len(db.ResourceKind) == 0 {
		return nil, errors.Internal().WithMessage("'ResourceKind' cannot be empty")
	}
	if len(db.ParentResourceID) == 0 {
		return nil, errors.Internal().WithMessage("'ParentResourceID' cannot be empty")
	}
	if len(db.ParentResourceKind) == 0 {
		return nil, errors.Internal().WithMessage("'ParentResourceKind' cannot be empty")
	}

	resourceKind, mapErr := mapTypeFromDB(db.ResourceKind)
	if mapErr != nil {
		return nil, mapErr
	}
	parentResourceKind, mapErr := mapTypeFromDB(db.ParentResourceKind)
	if mapErr != nil {
		return nil, mapErr
	}

	return &referentialintegrity.ReferentialIntegrityEntry{
		StandardResource: entities.StandardResource{
			StandardID: db.StandardID,
			Timestamps: entities.Timestamps{
				CreationDate: time.TimestampFromInt64(db.CreationDate),
				LastUpdate:   time.TimestampFromInt64(db.LastUpdate),
			},
		},
		ResourceID:         db.ResourceID,
		ResourceKind:       *resourceKind,
		ParentResourceID:   db.ParentResourceID,
		ParentResourceKind: *parentResourceKind,
	}, nil
}

func mapSliceFromDB(dbSlice []referentialintegritydb.ReferentialIntegrityEntryDB) ([]referentialintegrity.ReferentialIntegrityEntry, error) {
	items := make([]referentialintegrity.ReferentialIntegrityEntry, len(dbSlice))
	for index := range dbSlice {
		item, err := mapFromDB(dbSlice[index])
		if err != nil {
			return nil, err
		}
		items[index] = *item
	}

	return items, nil
}

func mapTypeToDB(resourceKind referentialintegrity.ResourceKind) (*string, error) {
	if resourceKind == referentialintegrity.KindAccount {
		k := referentialintegritydb.KindAccount
		return &k, nil
	}
	if resourceKind == referentialintegrity.KindApplication {
		k := referentialintegritydb.KindApplication
		return &k, nil
	}
	if resourceKind == referentialintegrity.KindHSMModule {
		k := referentialintegritydb.KindHSMModule
		return &k, nil
	}
	if resourceKind == referentialintegrity.KindHSMSlot {
		k := referentialintegritydb.KindHSMSlot
		return &k, nil
	}
	if resourceKind == referentialintegrity.KindUser {
		k := referentialintegritydb.KindUser
		return &k, nil
	}
	return nil, errors.Internal().WithMessage("couldn't map '%s' to a valid resource kind", resourceKind)
}

func mapTypeFromDB(resourceKind string) (*referentialintegrity.ResourceKind, error) {
	if resourceKind == referentialintegritydb.KindAccount {
		k := referentialintegrity.KindAccount
		return &k, nil
	}
	if resourceKind == referentialintegritydb.KindApplication {
		k := referentialintegrity.KindApplication
		return &k, nil
	}
	if resourceKind == referentialintegritydb.KindHSMModule {
		k := referentialintegrity.KindHSMModule
		return &k, nil
	}
	if resourceKind == referentialintegritydb.KindHSMSlot {
		k := referentialintegrity.KindHSMSlot
		return &k, nil
	}
	if resourceKind == referentialintegritydb.KindUser {
		k := referentialintegrity.KindUser
		return &k, nil
	}
	return nil, errors.Internal().WithMessage("couldn't map '%s' to a valid resource kind", resourceKind)
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
