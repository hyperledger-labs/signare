package accountdb

// AccountDB is the data struct of the resource in the database
type AccountDB struct {
	// InternalResourceID is the ID used to reference a resource internally in the application
	InternalResourceID string `storage:"internal_resource_id"`
	// Address is the Ethereum account to interact with the network
	Address string `storage:"address"`
	// ApplicationID the ID of the Application associated to this account
	ApplicationID string `storage:"application_id"`
	// UserID the ID of the User associated to this account
	UserID string `storage:"user_id"`
	// CreationDate is the timestamp of the moment of the creation of the resource
	CreationDate int64 `storage:"creation_date"`
	// LastUpdate is the timestamp of the moment of the last edition of the resource
	LastUpdate int64 `storage:"last_update"`
}

// AccountCreateDB is the data struct of the creation of a resource in the database
type AccountCreateDB struct {
	// AccountDB is the data struct of the resource in the database
	AccountDB
}

// AccountID is the primary key of an Account in the database
type AccountID struct {
	// Address is the Ethereum account to interact with the network
	Address string `storage:"address"`
	// ApplicationID the ID of the Application associated to this account
	ApplicationID string `storage:"application_id"`
	// UserID the ID of the User associated to this account
	UserID string `storage:"user_id"`
}

// AccountApplicationAddresses filters a list of accounts based on their Application
type AccountApplicationAddresses struct {
	// Address is the Ethereum account to interact with the network
	Address string `storage:"address"`
	// ApplicationID the ID of the Application associated to this account
	ApplicationID string `storage:"application_id"`
}
