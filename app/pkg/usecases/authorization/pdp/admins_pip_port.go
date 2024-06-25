package pdp

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// AdminsPolicyInformationPort is a port to adapt requests related to the Admin resource.
type AdminsPolicyInformationPort interface {
	// GetAdminRoles returns the list of roles assigned to an admin.
	GetAdminRoles(ctx context.Context, input GetAdminRolesInput) (*GetAdminRolesOutput, error)
}

// GetAdminRolesInput are the attributes needed to get the roles for an admin.
type GetAdminRolesInput struct {
	// AdminID is the ID of the user of the roles.
	AdminID entities.StandardID
}

// GetAdminRolesOutput is the result of getting the roles for an admin.
type GetAdminRolesOutput struct {
	// Roles list of roles assigned to a user.
	Roles []string
}
