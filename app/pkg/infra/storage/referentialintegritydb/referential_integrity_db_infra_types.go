package referentialintegritydb

import "github.com/hyperledger-labs/signare/app/pkg/entities"

const (
	KindAccount     = "account"
	KindApplication = "application"
	KindHSMModule   = "hardware_security_module"
	KindHSMSlot     = "hardware_security_module_slot"
	KindUser        = "user"
)

// ReferentialIntegrityEntryDB is the data struct of the resource in the database
type ReferentialIntegrityEntryDB struct {
	// StandardID is the ID of the resource
	entities.StandardID
	// ResourceID is the ID of the resource depending on the ParentResourceID
	ResourceID string `storage:"resource_id"`
	// ResourceKind is the kind of the resource depending on the ParentResourceID
	ResourceKind string `storage:"resource_kind"`
	// ResourceID is the ID of the parent resource
	ParentResourceID string `storage:"parent_resource_id"`
	// ParentResourceKind is the kind of the parent resource
	ParentResourceKind string `storage:"parent_resource_kind"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
}

// ReferentialIntegrityEntryCreateDB is the data struct of the creation of a resource in the database
type ReferentialIntegrityEntryCreateDB struct {
	// ReferentialIntegrityEntryDB is the data struct of the resource in the database
	ReferentialIntegrityEntryDB
}
