package pip

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
)

var _ pdp.AdminsPolicyInformationPort = (*DefaultAdminsPIPAdapter)(nil)

// GetAdminRoles returns the list of roles assigned to an admin
func (d DefaultAdminsPIPAdapter) GetAdminRoles(ctx context.Context, input pdp.GetAdminRolesInput) (*pdp.GetAdminRolesOutput, error) {
	getAdminRolesInput := admin.GetAdminInput{
		StandardID: input.AdminID,
	}
	getAdminOutput, getAdminErr := d.adminUseCase.GetAdmin(ctx, getAdminRolesInput)
	if getAdminErr != nil {
		return nil, getAdminErr
	}

	return &pdp.GetAdminRolesOutput{
		Roles: getAdminOutput.Roles,
	}, nil
}

// DefaultAdminsPIPAdapterOptions are the set of fields to create an DefaultAdminsPIPAdapter
type DefaultAdminsPIPAdapterOptions struct {
	// AdminUseCase defines the management of Admin in storage
	AdminUseCase admin.AdminUseCase
}

// DefaultAdminsPIPAdapter is a port to adapt requests related to the Admin resource
type DefaultAdminsPIPAdapter struct {
	adminUseCase admin.AdminUseCase
}

// ProvideDefaultAdminsPIPAdapter provides an instance of an DefaultAdminsPIPAdapter
func ProvideDefaultAdminsPIPAdapter(options DefaultAdminsPIPAdapterOptions) (*DefaultAdminsPIPAdapter, error) {
	return &DefaultAdminsPIPAdapter{
		adminUseCase: options.AdminUseCase,
	}, nil
}
