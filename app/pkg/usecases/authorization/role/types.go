package role

// Role is the name of the definition of the collection of permissions.
type Role struct {
	// ID is the name of the Role.
	ID string
}

// GetSupportedRolesInput are the attributes to fetch the supported collection of Role.
type GetSupportedRolesInput struct {
}

// GetSupportedRolesOutput is the supported collection of Role in the storage.
type GetSupportedRolesOutput struct {
	// Roles group of roles
	Roles []Role
}
