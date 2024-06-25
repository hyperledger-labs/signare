package role

import "context"

// RoleStorage defines the access methods to the Role storage.
type RoleStorage interface {
	// ListRoles fetches a collection of Role based on the ListRolesInput.
	ListRoles(ctx context.Context, input ListRolesInput) (*ListRolesOutput, error)
}

// ListRolesInput are the attributes to fetch the collection of Role.
type ListRolesInput struct {
}

// ListRolesOutput is the collection of Role in the storage.
type ListRolesOutput struct {
	// Roles is the array of Role.
	Roles []Role
}
