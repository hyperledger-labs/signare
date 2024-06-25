package roleinfile

// Role is the name of the definition of the collection of permissions
type Role struct {
	// ID is the name of the Role
	ID string `yaml:"id"`
	// Description describes the role with a phrase
	Description string `yaml:"description"`
	// Permissions are the array of permissions associated to the Role
	Permissions []string `yaml:"permissions"`
}

// RolesInfo is the data type that defines the roles in the file
type RolesInfo struct {
	// Default the default Role
	Default string `yaml:"default"`
	// Roles is the array of roles desribed in the file
	Roles []Role `yaml:"roles"`
}
