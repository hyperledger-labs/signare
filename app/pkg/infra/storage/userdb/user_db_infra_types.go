package userdb

import "github.com/hyperledger-labs/signare/app/pkg/entities"

// UserDB is the data struct of the resource in the database
type UserDB struct {
	// ApplicationStandardID is the ID of the resource
	entities.ApplicationStandardID
	// InternalResourceID is the ID used to reference a resource internally in the application
	InternalResourceID string `storage:"internal_resource_id"`
	// Roles are the roles in the RBAC system for the User
	Roles string `storage:"roles"`
	// Description of the resource
	Description string `storage:"description"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
	// ResourceVersion is the identifier of the current version of the resource
	ResourceVersion string `storage:"resource_version"`
}

// UserCreateDB is the data struct of the creation of a resource in the database
type UserCreateDB struct {
	// UserDB is the data struct of the resource in the database
	UserDB
}

// UserUpdateDB is the data struct of the update of a resource in the database
type UserUpdateDB struct {
	// UserDB is the data struct of the resource in the database
	UserDB
	// NewResourceVersion is the new resource version after the edition
	NewResourceVersion string `storage:"new_resource_version"`
}

// UserExistsDB is the data struct to check if a resource exists in the database
type UserExistsDB struct {
	// Exists is true if the resource exists
	Exists bool `storage:"exists_result" valid:"required"`
}
