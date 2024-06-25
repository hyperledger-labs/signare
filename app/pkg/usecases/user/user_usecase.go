package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/entities/address"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnection"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

// UserUseCase defines the management of the User resource.
type UserUseCase interface {
	// CreateUser creates a User. It returns the created User or an error if it fails.
	CreateUser(ctx context.Context, creation CreateUserInput) (*CreateUserOutput, error)
	// ListUsers returns all the Users or an error if it fails.
	ListUsers(ctx context.Context, listOptions ListUsersInput) (*ListUsersOutput, error)
	// GetUser returns the requested User or an error if it fails.
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	// EditUser edits a User. It returns the edited User or an error if it fails.
	EditUser(ctx context.Context, update EditUserInput) (*EditUserOutput, error)
	// DeleteUser deletes a User. It returns the deleted User or an error if it fails.
	DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error)
	// EnableAccounts adds accounts in the User's authorized accounts list. It returns the edited User or an error if it fails.
	EnableAccounts(ctx context.Context, input EnableAccountsInput) (*EnableAccountsOutput, error)
	// DisableAccount removes accounts in the User's authorized accounts list. It returns the edited User or an error if it fails.
	DisableAccount(ctx context.Context, input DisableAccountInput) (*DisableAccountOutput, error)
	AccountUseCase
}

func (u *DefaultUserUseCase) CreateUser(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input user")
	}

	if input.ID == nil {
		randomID := uuid.New().String()
		input.ID = &randomID
	}

	now := time.Now()
	user := User{
		ApplicationStandardResourceMeta: entities.ApplicationStandardResourceMeta{
			ApplicationStandardResource: entities.ApplicationStandardResource{
				ApplicationStandardID: entities.ApplicationStandardID{
					ID:            *input.ID,
					ApplicationID: input.ApplicationID,
				},
				Timestamps: entities.Timestamps{
					CreationDate: now,
					LastUpdate:   now,
				},
			},
		},
		Roles:       input.Roles,
		Description: input.Description,
	}
	user.InternalResourceID = entities.NewInternalResourceID()
	addUserToApplicationDependencyErr := u.addUserToApplicationDependency(ctx, user)
	if addUserToApplicationDependencyErr != nil {
		return nil, addUserToApplicationDependencyErr
	}

	getSupportedRolesInput := role.GetSupportedRolesInput{}
	getSupportedRolesOutput, getSupportedRolesErr := u.roleUseCase.GetSupportedRoles(ctx, getSupportedRolesInput)
	if getSupportedRolesErr != nil {
		return nil, getSupportedRolesErr
	}

	for _, role := range input.Roles {
		isSupported := false
		for _, supportedRole := range getSupportedRolesOutput.Roles {
			isSupported = supportedRole.ID == role
			if supportedRole.ID == role {
				break
			}
		}
		if !isSupported {
			return nil, errors.InvalidArgument().WithMessage("the role '%s' is not supported", role)
		}
	}

	addedUser, err := u.storage.Add(ctx, user)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err)
		}
		return nil, errors.InternalFromErr(err)
	}

	addedUser.Accounts = make([]Account, 0)
	return &CreateUserOutput{
		User: *addedUser,
	}, nil
}

func (u *DefaultUserUseCase) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.storage.Filter(input.ApplicationID)
	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	userCollection, err := u.storage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	for i, userItem := range userCollection.Items {
		user := userItem
		listAccountsInput := ListAccountsInput{
			UserID:        &user.ID,
			ApplicationID: user.ApplicationID,
		}
		userAccounts, errList := u.ListAccounts(ctx, listAccountsInput)
		if errList != nil {
			return nil, errors.InternalFromErr(errList)
		}
		accounts := make([]Account, len(userAccounts.Items))
		copy(accounts, userAccounts.Items)
		userCollection.Items[i].Accounts = accounts
	}
	return &ListUsersOutput{
		UserCollection: *userCollection,
	}, nil
}

func (u *DefaultUserUseCase) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	getInput := entities.ApplicationStandardID{
		ID:            input.ID,
		ApplicationID: input.ApplicationID,
	}
	user, err := u.storage.Get(ctx, getInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage("user [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	listAccountsInput := ListAccountsInput{
		ApplicationID: user.ApplicationID,
		UserID:        &user.ID,
	}
	userAccounts, errList := u.ListAccounts(ctx, listAccountsInput)
	if errList != nil {
		return nil, errors.InternalFromErr(errList)
	}

	user.Accounts = userAccounts.Items
	return &GetUserOutput{
		User: *user,
	}, nil
}

func (u *DefaultUserUseCase) EditUser(ctx context.Context, input EditUserInput) (*EditUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	if len(input.Roles) < 1 {
		return nil, errors.InvalidArgument().WithMessage("the 'Role' cannot be empty")
	}

	getSupportedRolesInput := role.GetSupportedRolesInput{}
	getSupportedRolesOutput, getSupportedRolesErr := u.roleUseCase.GetSupportedRoles(ctx, getSupportedRolesInput)
	if getSupportedRolesErr != nil {
		return nil, getSupportedRolesErr
	}

	for _, role := range input.Roles {
		isSupported := false
		for _, supportedRole := range getSupportedRolesOutput.Roles {
			isSupported = supportedRole.ID == role
			if supportedRole.ID == role {
				break
			}
		}
		if !isSupported {
			msg := fmt.Sprintf("the role '%s' is not supported", role)
			return nil, errors.InvalidArgument().WithMessage(msg).SetHumanReadableMessage(msg)
		}
	}

	user := User{
		ApplicationStandardResourceMeta: entities.ApplicationStandardResourceMeta{
			ApplicationStandardResource: entities.ApplicationStandardResource{
				ApplicationStandardID: entities.ApplicationStandardID{
					ID:            input.ID,
					ApplicationID: input.ApplicationID,
				},
				Timestamps: entities.Timestamps{
					LastUpdate: time.Now(),
				},
			},
			ResourceVersion: input.ResourceVersion,
		},
		Roles: input.Roles,
	}
	if input.Description != nil {
		user.Description = input.Description
	}

	editedUser, err := u.storage.Edit(ctx, user)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage("user [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	listAccountsInput := ListAccountsInput{
		ApplicationID: user.ApplicationID,
		UserID:        &user.ID,
	}
	userAccounts, errList := u.ListAccounts(ctx, listAccountsInput)
	if errList != nil {
		return nil, errors.InternalFromErr(errList)
	}

	editedUser.Accounts = userAccounts.Items

	return &EditUserOutput{
		User: *editedUser,
	}, nil
}

func (u *DefaultUserUseCase) DeleteUser(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}
	removeAllDependenciesErr := u.removeAllUserDependencies(ctx, input.ApplicationStandardID)
	if removeAllDependenciesErr != nil {
		return nil, removeAllDependenciesErr
	}

	removeInput := entities.ApplicationStandardID{
		ID:            input.ID,
		ApplicationID: input.ApplicationID,
	}
	user, err := u.storage.Remove(ctx, removeInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage("user [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteUserOutput{
		User: *user,
	}, nil
}

func (u *DefaultUserUseCase) EnableAccounts(ctx context.Context, input EnableAccountsInput) (*EnableAccountsOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	tracer := logger.NewTracer(ctx)
	tracer.AddProperty("user", input.UserID)
	tracer.AddProperty("application", input.ApplicationID)
	tracer.AddProperty("addresses", input.Addresses)
	tracer.Debug("enabling accounts of user")

	getUserInput := entities.ApplicationStandardID{
		ID:            input.UserID,
		ApplicationID: input.ApplicationID,
	}
	user, err := u.storage.Get(ctx, getUserInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("could not find user [%s] on application [%s]", input.UserID, input.ApplicationID)
		}
		return nil, errors.InternalFromErr(err)
	}

	byApplicationInput := hsmconnection.ByApplicationInput{
		ApplicationID: input.ApplicationID,
	}

	hsmConnection, byApplicationErr := u.hsmConnectionResolver.ByApplication(ctx, byApplicationInput)
	if byApplicationErr != nil {
		return nil, byApplicationErr
	}

	// Accounts need to be validated with the HSM manager to see if they exist in their slots.
	listAddressesInput := hsmconnector.ListAddressesInput{
		SlotConnectionData: hsmconnector.SlotConnectionData{
			Slot:       hsmConnection.Slot,
			Pin:        hsmConnection.Pin,
			ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			ChainID:    hsmConnection.ChainID,
		},
	}
	listAddressesOutput, err := u.hsmConnector.ListAddresses(ctx, listAddressesInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	if !areAddressesValid(listAddressesOutput.Items, input.Addresses) {
		msg := fmt.Sprintf("one or more accounts '%s' do not exist in the HSM", input.Addresses)
		return nil, errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	accountsToCreate := make([]CreateAccountInput, len(input.Addresses))
	for i, addr := range input.Addresses {
		a := CreateAccountInput{
			AccountID: AccountID{
				Address:       addr,
				UserID:        input.UserID,
				ApplicationID: input.ApplicationID,
			},
		}
		accountsToCreate[i] = a
	}

	var accountsNotAdded = make([]string, 0)
	for _, createAccountInput := range accountsToCreate {
		_, createAccountErr := u.CreateAccount(ctx, createAccountInput)
		if createAccountErr != nil && !errors.IsAlreadyExists(createAccountErr) {
			accountsNotAdded = append(accountsNotAdded, createAccountInput.Address.String())
			continue
		}
	}

	if len(accountsNotAdded) > 0 {
		formattedAddresses := fmt.Sprintf("[%s]", strings.Join(accountsNotAdded, ", "))
		return nil, errors.Internal().WithMessage("error while adding accounts to the user. The following accounts were not added: %s", formattedAddresses)
	}

	listAccountsInput := ListAccountsInput{
		ApplicationID: input.ApplicationID,
		UserID:        &input.UserID,
	}
	userAccounts, err := u.ListAccounts(ctx, listAccountsInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	user.Accounts = userAccounts.Items
	return &EnableAccountsOutput{
		User: *user,
	}, nil
}

func (u *DefaultUserUseCase) DisableAccount(ctx context.Context, input DisableAccountInput) (*DisableAccountOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	tracer := logger.NewTracer(ctx)
	tracer.AddProperty("user", input.UserID)
	tracer.AddProperty("application", input.ApplicationID)
	tracer.AddProperty("address", input.Address.String())
	tracer.Debug("disabling account of user")

	deleteAccountInput := DeleteAccountInput{
		AccountID: AccountID{
			Address:       input.Address,
			UserID:        input.UserID,
			ApplicationID: input.ApplicationID,
		},
	}
	_, err = u.DeleteAccount(ctx, deleteAccountInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("account [%s] not found for user [%s]", input.Address, input.UserID)
		}
		return nil, errors.InternalFromErr(err)
	}

	getUserInput := entities.ApplicationStandardID{
		ID:            input.UserID,
		ApplicationID: input.ApplicationID,
	}
	user, err := u.storage.Get(ctx, getUserInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("could not find user [%s] on application [%s]", input.UserID, input.ApplicationID)
		}
		return nil, errors.InternalFromErr(err)
	}

	listAccountsInput := ListAccountsInput{
		ApplicationID: input.ApplicationID,
		UserID:        &input.UserID,
	}
	userAccounts, err := u.ListAccounts(ctx, listAccountsInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	user.Accounts = userAccounts.Items
	return &DisableAccountOutput{
		User: *user,
	}, nil
}

func areAddressesValid(items []address.Address, addresses []address.Address) bool {
	for _, addr := range addresses {
		if !containsAddress(items, addr) {
			return false
		}
	}
	return true
}

func containsAddress(items []address.Address, address address.Address) bool {
	for _, item := range items {
		if item == address {
			return true
		}
	}
	return false
}

var _ UserUseCase = new(DefaultUserUseCase)

// DefaultUserUseCase default management of User in configuration implementation.
type DefaultUserUseCase struct {
	// storage is the persistence adapter of the User.
	storage UserStorage
	// accountStorage is the persistence adapter of the Account.
	accountStorage AccountStorage

	// applicationUseCase defines how to interact with Application resources.
	applicationUseCase application.ApplicationUseCase
	// hsmConnectionResolver finds what HSMConnection is required depending on the constraints.
	hsmConnectionResolver hsmconnection.Resolver
	// hsmConnector connects with the HSM and operates with it.
	hsmConnector hsmconnector.HSMConnector
	// referentialIntegrityUseCase to manage dependencies between resources.
	referentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
	// roleUseCase defines how to interact with Role resources.
	roleUseCase role.RoleUseCase
}

// DefaultUserUseCaseOptions configures a DefaultUserUseCase.
type DefaultUserUseCaseOptions struct {
	// Storage is the persistence adapter of the User.
	Storage UserStorage
	// AccountStorage is the persistence adapter of the Account.
	AccountStorage AccountStorage

	// ApplicationUseCase defines how to interact with Application resources.
	ApplicationUseCase application.ApplicationUseCase
	// HSMConnectionResolver finds what HSMConnection is required depending on the constraints.
	HSMConnectionResolver hsmconnection.Resolver
	// HSMConnector connects with the HSM and operates with it.
	HSMConnector hsmconnector.HSMConnector
	// ReferentialIntegrityUseCase to manage dependencies between resources.
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
	// RoleUseCase defines how to interact with Role resources.
	RoleUseCase role.RoleUseCase
}

// ProvideDefaultUseCase creates a DefaultUserUseCase with the given options.
func ProvideDefaultUseCase(options DefaultUserUseCaseOptions) (*DefaultUserUseCase, error) {
	if options.Storage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'Storage' was not provided")
	}
	if options.AccountStorage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'AccountStorage' was not provided")
	}
	if options.ApplicationUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ApplicationUseCase' was not provided")
	}
	if options.HSMConnectionResolver == nil {
		return nil, errors.Internal().WithMessage("mandatory 'Resolver' was not provided")
	}
	if options.HSMConnector == nil {
		return nil, errors.Internal().WithMessage("mandatory 'HSMConnector' was not provided")
	}
	if options.ReferentialIntegrityUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ReferentialIntegrityUseCase' was not provided")
	}
	if options.RoleUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'RoleUseCase' was not provided")
	}
	return &DefaultUserUseCase{
		storage:                     options.Storage,
		applicationUseCase:          options.ApplicationUseCase,
		accountStorage:              options.AccountStorage,
		roleUseCase:                 options.RoleUseCase,
		hsmConnector:                options.HSMConnector,
		hsmConnectionResolver:       options.HSMConnectionResolver,
		referentialIntegrityUseCase: options.ReferentialIntegrityUseCase,
	}, nil
}
