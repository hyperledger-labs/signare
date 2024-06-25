package applicationdb

import "github.com/hyperledger-labs/signare/app/pkg/entities"

// ApplicationDB is the data struct of the resource in the database
type ApplicationDB struct {
	// StandardID is the ID of the resource
	entities.StandardID
	// InternalResourceID is the ID used to reference a resource internally in the application
	InternalResourceID string `storage:"internal_resource_id"`
	// ChainID is the id of the Ethereum network in which the application operates
	ChainID string `storage:"chain_id"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
	// ResourceVersion is the identifier of the current version of the resource
	ResourceVersion string `storage:"resource_version"`
	// Description of the resource
	Description *string `storage:"description"`
}

// ApplicationCreateDB is the data struct of the creation of a resource in the database
type ApplicationCreateDB struct {
	// ApplicationDB is the data struct of the resource in the database
	ApplicationDB
}

// ApplicationUpdateDB is the data struct of the update of a resource in the database
type ApplicationUpdateDB struct {
	// ApplicationDB is the data struct of the resource in the database
	ApplicationDB
	// NewResourceVersion is the new resource version after the edition
	NewResourceVersion string `storage:"new_resource_version"`
}

// ApplicationExistsDB is the data struct to check if a resource exists in the database
type ApplicationExistsDB struct {
	// Exists is true if the resource exists
	Exists bool `storage:"exists_result" valid:"required"`
}
