package hsmmoduledb

import "github.com/hyperledger-labs/signare/app/pkg/entities"

const (
	SoftHSMModuleKind = "SoftHSM"
)

// HardwareSecurityModuleDB is the data struct of the resource in the database
type HardwareSecurityModuleDB struct {
	// StandardID is the ID of the resource
	entities.StandardID
	// InternalResourceID is the ID used to reference a resource internally in the application
	InternalResourceID string `storage:"internal_resource_id"`
	// Kind of the HSM
	Kind string `storage:"kind"`
	// Configuration of the HSM
	Configuration string `storage:"configuration"`
	// Description of the resource
	Description string `storage:"description"`
	// ResourceVersion is the identifier of the current version of the resource
	ResourceVersion string `storage:"resource_version"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
}

// HardwareSecurityModuleCreateDB is the data struct of the creation of a resource in the database
type HardwareSecurityModuleCreateDB struct {
	// HardwareSecurityModuleDB is the data struct of the resource in the database
	HardwareSecurityModuleDB
}

// HardwareSecurityModuleUpdateDB is the data struct of the update of a resource in the database
type HardwareSecurityModuleUpdateDB struct {
	// HardwareSecurityModuleDB is the data struct of the resource in the database
	HardwareSecurityModuleDB
	// NewResourceVersion is the new resource version after the edition
	NewResourceVersion string `storage:"new_resource_version"`
}

// HardwareSecurityModuleExistsDB is the data struct to check if a resource exists in the database
type HardwareSecurityModuleExistsDB struct {
	// Exists is true if the resource exists
	Exists bool `storage:"exists_result" valid:"required"`
}
