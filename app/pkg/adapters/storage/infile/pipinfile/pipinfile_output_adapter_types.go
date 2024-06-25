package pipinfile

// Role defines a role to use for access controls.
type Role struct {
	// ID identifier of the role.
	ID string `yaml:"id"`
	// Description role description.
	Description string `yaml:"description"`
	// Permissions permissions allowed for the role.
	Permissions []string `yaml:"permissions"`
}

// RolesInfo information about roles.
type RolesInfo struct {
	// Roles group of roles.
	Roles []Role `yaml:"roles"`
}

// Permission contains a set of actions.
type Permission struct {
	// ID identifier of the permission.
	ID string `yaml:"id"`
	// Description of the permission.
	Description string `yaml:"description"`
	// Actions group of actions the permission grants access to.
	Actions []string `yaml:"actions"`
}

// PermissionCollection group of Permission.
type PermissionCollection struct {
	Permissions []Permission `yaml:"permissions"`
}

// ActionCollection group of actions.
type ActionCollection struct {
	Actions []string `yaml:"actions"`
}

// mergeActions concatenates returns an ActionCollection after concatenating actions from the two inputs
func mergeActions(actionsOne, actionsTwo ActionCollection) ActionCollection {
	return ActionCollection{
		Actions: append(actionsOne.Actions, actionsTwo.Actions...),
	}
}
