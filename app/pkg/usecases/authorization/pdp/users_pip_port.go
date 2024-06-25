package pdp

import "context"

// UsersPolicyInformationPort is a port to adapt requests related to the User resource.
type UsersPolicyInformationPort interface {
	// GetUserRoles returns the list of roles assigned to a user.
	GetUserRoles(ctx context.Context, input GetUserRolesInput) (*GetUserRolesOutput, error)
}

// GetUserRolesInput are the attributes needed to get the roles for a user.
type GetUserRolesInput struct {
	// UserID is the ID of the user.
	UserID string
	// ApplicationID is the ID of the application.
	ApplicationID string
}

// GetUserRolesOutput is the result of getting the roles for a user.
type GetUserRolesOutput struct {
	// Roles list of roles assigned to a user.
	Roles []string
}
