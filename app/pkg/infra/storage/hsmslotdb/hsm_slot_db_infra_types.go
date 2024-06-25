package hsmslotdb

import "github.com/hyperledger-labs/signare/app/pkg/entities"

// HSMSlotDB is the data struct of the resource in the database
type HSMSlotDB struct {
	// StandardID is the ID of the resource
	entities.StandardID
	// InternalResourceID is the ID used to reference a resource internally in the application
	InternalResourceID string `storage:"internal_resource_id"`
	// ApplicationID the ID of the Application associated to this HSM Slot
	ApplicationID string `storage:"application_id"`
	// HSMModuleID the ID of the HSM associated to this HSM Slot
	HSMModuleID string `storage:"hardware_security_module_id"`
	// Slot identifier within the HSM
	Slot string `storage:"slot"`
	// Pin the password of the HSM Slot in the HSM
	Pin string `storage:"pin"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
	// ResourceVersion is the identifier of the current version of the resource
	ResourceVersion string `storage:"resource_version"`
}

// HSMSlotCreateDB is the data struct of the resource in the database
type HSMSlotCreateDB struct {
	// HSMSlotDB is the data struct of the resource in the database
	HSMSlotDB
}

// HSMSlotUpdatePinDB is the data struct of the update of a resource in the database
type HSMSlotUpdatePinDB struct {
	// StandardID is the ID of the resource
	entities.StandardID
	// ResourceVersion is the identifier of the current version of the resource
	ResourceVersion string `storage:"resource_version"`
	// Pin the password of the HSM Slot in the HSM
	Pin string `storage:"pin"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
	// NewResourceVersion is the new resource version after the edition
	NewResourceVersion string `storage:"new_resource_version"`
}

// HSMSlotExistsDB is the data struct to check if a resource exists in the database
type HSMSlotExistsDB struct {
	// Exists is true if the resource exists
	Exists bool `storage:"exists_result" valid:"required"`
}
