package user

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"
)

// CreateAccount implements AccountUseCase's CreateAccount to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) CreateAccount(ctx context.Context, input CreateAccountInput) (*CreateAccountOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.createAccountInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*CreateAccountOutput), nil
}

// ListAccounts implements AccountUseCase's ListAccounts to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) ListAccounts(ctx context.Context, input ListAccountsInput) (*ListAccountsOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.listAccountsInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*ListAccountsOutput), nil
}

// GetAccount implements AccountUseCase's GetAccount to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) GetAccount(ctx context.Context, input GetAccountInput) (*GetAccountOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.getAccountInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*GetAccountOutput), nil
}

// DeleteAccount implements AccountUseCase's DeleteAccount to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) DeleteAccount(ctx context.Context, input DeleteAccountInput) (*DeleteAccountOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.deleteAccountInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*DeleteAccountOutput), nil
}

// DeleteAllAccountsForAddress implements AccountUseCase's DeleteAllAccountsForAddress to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) DeleteAllAccountsForAddress(ctx context.Context, input DeleteAllAccountsForAddressInput) (*DeleteAllAccountsForAddressOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.deleteAllAccountsForAddressInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*DeleteAllAccountsForAddressOutput), nil
}

var _ AccountUseCase = new(DefaultUseCaseTransactionalDecorator)

// DefaultUseCaseTransactionalDecorator decorates struct AccountUseCase wrapped with a transactional manager
type DefaultUseCaseTransactionalDecorator struct {
	// AccountUseCase is the usecase to be decorated
	AccountUseCase
	// transactionalManager defines the functionality to execute a transaction in a transactional manner.
	transactionalManager transactionalmanager.TransactionalManagerUseCase
}

// DefaultUseCaseTransactionalDecoratorOptions is the structure representing the DefaultUseCaseTransactionalDecorator dependencies.
type DefaultUseCaseTransactionalDecoratorOptions struct {
	// AccountUseCase is the usecase to be decorated
	AccountUseCase AccountUseCase
	// TransactionalManager defines the functionality to execute a transaction in a transactional manner.
	TransactionalManager transactionalmanager.TransactionalManagerUseCase
}

// ProvideDefaultUseCaseTransactionalDecorator creates a DefaultUseCaseTransactionalDecorator with the given options.
func ProvideDefaultUseCaseTransactionalDecorator(options DefaultUseCaseTransactionalDecoratorOptions) (*DefaultUseCaseTransactionalDecorator, error) {
	if options.AccountUseCase == nil {
		errorMessage := "'AccountUseCase' is mandatory"
		return nil, errors.InvalidArgument().WithMessage(errorMessage)
	}
	if options.TransactionalManager == nil {
		errorMessage := "'TransactionalManager' is mandatory"
		return nil, errors.InvalidArgument().WithMessage(errorMessage)
	}
	return &DefaultUseCaseTransactionalDecorator{
		AccountUseCase:       options.AccountUseCase,
		transactionalManager: options.TransactionalManager,
	}, nil
}

func (_d *DefaultUseCaseTransactionalDecorator) createAccountInternal(_ context.Context, input CreateAccountInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.AccountUseCase.CreateAccount(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) listAccountsInternal(_ context.Context, generationalManagedContractListOptions ListAccountsInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.AccountUseCase.ListAccounts(ctx2, generationalManagedContractListOptions)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) getAccountInternal(_ context.Context, input GetAccountInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.AccountUseCase.GetAccount(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) deleteAccountInternal(_ context.Context, input DeleteAccountInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.AccountUseCase.DeleteAccount(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) deleteAllAccountsForAddressInternal(_ context.Context, input DeleteAllAccountsForAddressInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.AccountUseCase.DeleteAllAccountsForAddress(ctx2, input)
	}
}
