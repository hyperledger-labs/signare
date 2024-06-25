package user

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnection"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"

	"github.com/asaskevich/govalidator"
)

// AccountUseCase defines the management of the Account resource.
type AccountUseCase interface {
	// CreateAccount creates an Account. It returns the created Account or an error if it fails.
	CreateAccount(ctx context.Context, input CreateAccountInput) (*CreateAccountOutput, error)
	// ListAccounts returns all the Accounts or an error if it fails.
	ListAccounts(ctx context.Context, input ListAccountsInput) (*ListAccountsOutput, error)
	// GetAccount returns the requested Account or an error if it fails.
	GetAccount(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error)
	// DeleteAccount deletes an Account. It returns the deleted Account or an error if it fails.
	DeleteAccount(ctx context.Context, input DeleteAccountInput) (*DeleteAccountOutput, error)
	// DeleteAllAccountsForAddress deletes all the Accounts for a given application with a specific address removing this address from the HSM. It returns an error if it fails.
	DeleteAllAccountsForAddress(ctx context.Context, input DeleteAllAccountsForAddressInput) (*DeleteAllAccountsForAddressOutput, error)
}

func (u *DefaultUserUseCase) CreateAccount(ctx context.Context, input CreateAccountInput) (*CreateAccountOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		logger.LogEntry(ctx).Debugf("couldn't validate input account: %s", err.Error())
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input account")
	}
	account := Account{
		AccountID: AccountID{
			Address:       input.Address,
			UserID:        input.UserID,
			ApplicationID: input.ApplicationID,
		},
		Timestamps: entities.Timestamps{
			CreationDate: time.Now(),
			LastUpdate:   time.Now(),
		},
	}
	account.InternalResourceID = entities.NewInternalResourceID()
	addAccountToUserDependencyErr := u.addAccountToUserDependency(ctx, account)
	if addAccountToUserDependencyErr != nil {
		return nil, addAccountToUserDependencyErr
	}

	addedAccount, err := u.accountStorage.Add(ctx, account)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err).SetHumanReadableMessage("addedAccount [%s] already exists", input.AccountID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &CreateAccountOutput{
		Account: *addedAccount,
	}, nil
}

func (u *DefaultUserUseCase) ListAccounts(ctx context.Context, input ListAccountsInput) (*ListAccountsOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.accountStorage.Filter(input.ApplicationID)
	if input.UserID != nil {
		filters.FilterByUserID(*input.UserID)
	}
	if input.Address != nil {
		filters.FilterByAddress(input.Address.String())
	}

	collection, err := u.accountStorage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListAccountsOutput{
		AccountCollection: *collection,
	}, nil
}

func (u *DefaultUserUseCase) GetAccount(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	account, err := u.accountStorage.Get(ctx, input.AccountID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("account [%s] not found", input.AccountID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &GetAccountOutput{
		Account: *account,
	}, nil
}

func (u *DefaultUserUseCase) DeleteAccount(ctx context.Context, input DeleteAccountInput) (*DeleteAccountOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	removeAllAccountDependencies := u.removeAllAccountDependencies(ctx, input.AccountID)
	if removeAllAccountDependencies != nil {
		return nil, removeAllAccountDependencies
	}

	account, err := u.accountStorage.Remove(ctx, input.AccountID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("account [%s] not found", input.AccountID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteAccountOutput{
		Account: *account,
	}, nil
}

func (u *DefaultUserUseCase) DeleteAllAccountsForAddress(ctx context.Context, input DeleteAllAccountsForAddressInput) (*DeleteAllAccountsForAddressOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	byApplicationInput := hsmconnection.ByApplicationInput{
		ApplicationID: input.ApplicationID,
	}
	hsmConnection, byApplicationErr := u.hsmConnectionResolver.ByApplication(ctx, byApplicationInput)
	if byApplicationErr != nil {
		return nil, byApplicationErr
	}

	tracer := logger.NewTracer(ctx)
	tracer.AddProperty("address", input.Address.String())
	tracer.AddProperty("moduleKind", hsmConnection.ModuleKind)
	tracer.AddProperty("slot", hsmConnection.Slot)
	tracer.Debug("removing address from HSM")

	// 1. Remove it from the HSM
	removeAddressInput := hsmconnector.RemoveAddressInput{
		SlotConnectionData: hsmconnector.SlotConnectionData{
			Slot:       hsmConnection.Slot,
			Pin:        hsmConnection.Pin,
			ModuleKind: hsmconnector.ModuleKind(hsmConnection.ModuleKind),
			ChainID:    hsmConnection.ChainID,
		},
		Address: input.Address,
	}
	_, err = u.hsmConnector.RemoveAddress(ctx, removeAddressInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err)
		}
		return nil, errors.InternalFromErr(err)
	}

	tracer.Trace("removed address from HSM")

	// 2. Remove it from storage
	accounts, err := u.accountStorage.RemoveAllForAddress(ctx, input.ApplicationID, input.Address)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteAllAccountsForAddressOutput{
		Items: accounts.Items,
	}, nil
}

var _ AccountUseCase = new(DefaultUserUseCase)
