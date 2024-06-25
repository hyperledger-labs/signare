package pdp

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"

	"github.com/asaskevich/govalidator"
)

// PolicyDecisionPointUseCase is the business logic to perform user authorization for different actions.
type PolicyDecisionPointUseCase interface {
	// AuthorizeUserAccount checks if a user is authorized to use an account, returns an error if it doesn't.
	AuthorizeUserAccount(ctx context.Context, input AuthorizeUserAccountInput) (*AuthorizeUserAccountOutput, error)
	// AuthorizeUser checks if the user is authorized to perform an action, returns an error if it doesn't.
	AuthorizeUser(ctx context.Context, input AuthorizeUserInput) (*AuthorizeUserOutput, error)
}

// AuthorizeUserAccount checks if a user is authorized to use an account, returns an error if it doesn't.
func (useCase DefaultPolicyDecisionPointUseCase) AuthorizeUserAccount(ctx context.Context, input AuthorizeUserAccountInput) (*AuthorizeUserAccountOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}
	getAccountInput := GetAccountInput{
		AccountID: AccountID{
			UserID:        input.UserID,
			ApplicationID: input.ApplicationID,
			Address:       input.Address,
		},
	}
	_, getAccountErr := useCase.accountsPolicyInformationAdapter.GetAccount(ctx, getAccountInput)
	if getAccountErr != nil {
		return nil, getAccountErr
	}

	return &AuthorizeUserAccountOutput{}, nil
}

// AuthorizeUser checks if the user is authorized to perform an action, returns an error if it's not.
func (useCase DefaultPolicyDecisionPointUseCase) AuthorizeUser(ctx context.Context, input AuthorizeUserInput) (*AuthorizeUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	var roles *[]string
	if input.ApplicationID != nil && *input.ApplicationID != "" {
		getUserInput := GetUserRolesInput{
			UserID:        input.UserID,
			ApplicationID: *input.ApplicationID,
		}
		getUserRolesOutput, getUserRolesErr := useCase.usersPolicyInformationAdapter.GetUserRoles(ctx, getUserInput)
		if getUserRolesErr != nil {
			return nil, getUserRolesErr
		}

		if getUserRolesOutput != nil {
			roles = &getUserRolesOutput.Roles
		}
	}

	if roles == nil {
		getAdminInput := GetAdminRolesInput{
			AdminID: entities.StandardID{
				ID: input.UserID,
			},
		}
		getAdminRolesOutput, getAdminRolesErr := useCase.adminsPolicyInformationAdapter.GetAdminRoles(ctx, getAdminInput)
		if getAdminRolesErr != nil {
			return nil, getAdminRolesErr
		}
		roles = &getAdminRolesOutput.Roles
	}

	listAllowedActionsInput := ListActionsInput{
		Roles: *roles,
	}
	listActionsOutput, listAllowedActionsErr := useCase.actionsPolicyInformationPointPort.ListActions(ctx, listAllowedActionsInput)
	if listAllowedActionsErr != nil {
		return nil, listAllowedActionsErr
	}

	if _, ok := listActionsOutput.Actions.actionSet[input.ActionID]; !ok {
		return nil, errors.PreconditionFailed().SetHumanReadableMessage("action not authorized for user [%s]", input.UserID)
	}

	return &AuthorizeUserOutput{}, nil
}

var _ PolicyDecisionPointUseCase = (*DefaultPolicyDecisionPointUseCase)(nil)

// DefaultPolicyDecisionPointUseCaseOptions are the set of fields to create an DefaultPolicyDecisionPointUseCase
type DefaultPolicyDecisionPointUseCaseOptions struct {
	// AccountsPolicyInformationPort is a port to adapt requests related to the Account resource
	AccountsPolicyInformationAdapter AccountsPolicyInformationPort
	// ActionsPolicyInformationPointPort is a port to adapt requests related to the Actions.
	ActionsPolicyInformationPointPort ActionsPolicyInformationPointPort
	// AdminsPolicyInformationAdapter is a port to adapt requests related to the Admin resource.
	AdminsPolicyInformationAdapter AdminsPolicyInformationPort
	// UsersPolicyInformationAdapter is a port to adapt requests related to the User resource.
	UsersPolicyInformationAdapter UsersPolicyInformationPort
}

// DefaultPolicyDecisionPointUseCase is the business logic to perform user authorization for different actions.
type DefaultPolicyDecisionPointUseCase struct {
	accountsPolicyInformationAdapter  AccountsPolicyInformationPort
	actionsPolicyInformationPointPort ActionsPolicyInformationPointPort
	adminsPolicyInformationAdapter    AdminsPolicyInformationPort
	usersPolicyInformationAdapter     UsersPolicyInformationPort
}

// ProvideDefaultPolicyDecisionPointUseCase provides an instance of an DefaultPolicyDecisionPointUseCase.
func ProvideDefaultPolicyDecisionPointUseCase(options DefaultPolicyDecisionPointUseCaseOptions) (*DefaultPolicyDecisionPointUseCase, error) {
	if options.UsersPolicyInformationAdapter == nil {
		return nil, errors.Internal().WithMessage("mandatory 'UsersPolicyInformationAdapter' not provided")
	}
	if options.AccountsPolicyInformationAdapter == nil {
		return nil, errors.Internal().WithMessage("mandatory 'AccountsPolicyInformationAdapter' not provided")
	}
	if options.ActionsPolicyInformationPointPort == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ActionsPolicyInformationPointPort' not provided")
	}
	if options.AdminsPolicyInformationAdapter == nil {
		return nil, errors.Internal().WithMessage("mandatory 'AdminsPolicyInformationAdapter' not provided")
	}
	if options.UsersPolicyInformationAdapter == nil {
		return nil, errors.Internal().WithMessage("mandatory 'UsersPolicyInformationAdapter' not provided")
	}

	return &DefaultPolicyDecisionPointUseCase{
		accountsPolicyInformationAdapter:  options.AccountsPolicyInformationAdapter,
		actionsPolicyInformationPointPort: options.ActionsPolicyInformationPointPort,
		adminsPolicyInformationAdapter:    options.AdminsPolicyInformationAdapter,
		usersPolicyInformationAdapter:     options.UsersPolicyInformationAdapter,
	}, nil
}
