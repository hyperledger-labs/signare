package pip

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/pdp"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

var _ pdp.UsersPolicyInformationPort = (*DefaultUsersPIPAdapter)(nil)

// GetUserRoles returns the list of roles assigned to a user
func (d DefaultUsersPIPAdapter) GetUserRoles(ctx context.Context, input pdp.GetUserRolesInput) (*pdp.GetUserRolesOutput, error) {
	getUserInput := user.GetUserInput{
		ApplicationStandardID: entities.ApplicationStandardID{
			ID:            input.UserID,
			ApplicationID: input.ApplicationID,
		},
	}

	getUserOutput, getUserErr := d.userUseCase.GetUser(ctx, getUserInput)
	if getUserErr != nil {
		return nil, getUserErr
	}

	return &pdp.GetUserRolesOutput{
		Roles: getUserOutput.Roles,
	}, nil
}

// DefaultUsersPIPAdapterOptions are the set of fields to create an DefaultUsersPIPAdapter
type DefaultUsersPIPAdapterOptions struct {
	// UserUseCase defines the management of the User resource
	UserUseCase user.UserUseCase
}

// DefaultUsersPIPAdapter is a port to adapt requests related to the User resource
type DefaultUsersPIPAdapter struct {
	userUseCase user.UserUseCase
}

// ProvideDefaultUsersPIPAdapter provides an instance of an DefaultUsersPIPAdapter
func ProvideDefaultUsersPIPAdapter(options DefaultUsersPIPAdapterOptions) (*DefaultUsersPIPAdapter, error) {
	return &DefaultUsersPIPAdapter{
		userUseCase: options.UserUseCase,
	}, nil
}
