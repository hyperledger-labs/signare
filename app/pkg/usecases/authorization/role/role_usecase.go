package role

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
)

// RoleUseCase the set of use cases for the Roles.
type RoleUseCase interface {
	// GetSupportedRoles fetches the list of supported roles for the signare users.
	GetSupportedRoles(ctx context.Context, input GetSupportedRolesInput) (*GetSupportedRolesOutput, error)
}

// GetSupportedRoles fetches the list of supported roles for the signare users.
func (d DefaultRoleUseCase) GetSupportedRoles(ctx context.Context, _ GetSupportedRolesInput) (*GetSupportedRolesOutput, error) {
	listRolesInput := ListRolesInput{}
	listRolesOutput, listRolesErr := d.roleStorage.ListRoles(ctx, listRolesInput)
	if listRolesErr != nil {
		return nil, listRolesErr
	}

	return &GetSupportedRolesOutput{
		Roles: listRolesOutput.Roles,
	}, nil
}

var _ RoleUseCase = new(DefaultRoleUseCase)

// DefaultRoleUseCaseOptions are the set of fields to create an DefaultRoleUseCase.
type DefaultRoleUseCaseOptions struct {
	// RoleStorage defines the access methods to the Role storage
	RoleStorage RoleStorage
}

// DefaultRoleUseCase is the struct that implements the RoleUseCase for the RBAC.
type DefaultRoleUseCase struct {
	roleStorage RoleStorage
}

// ProvideDefaultRoleUseCase provides an instance of an DefaultRoleUseCase.
func ProvideDefaultRoleUseCase(options DefaultRoleUseCaseOptions) (*DefaultRoleUseCase, error) {
	if options.RoleStorage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'RoleStorage' not provided")
	}

	return &DefaultRoleUseCase{
		roleStorage: options.RoleStorage,
	}, nil
}
