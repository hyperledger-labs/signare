package hsmmodule

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"
)

// CreateHSMModule implements DefaultUseCase's CreateHSMModule to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) CreateHSMModule(ctx context.Context, input CreateHSMModuleInput) (*CreateHSMModuleOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.createHSMModuleInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*CreateHSMModuleOutput), nil
}

// ListHSMModules implements DefaultUseCase's ListHSMModules to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) ListHSMModules(ctx context.Context, input ListHSMModulesInput) (*ListHSMModulesOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.listHSMModulesInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*ListHSMModulesOutput), nil
}

// GetHSMModule implements DefaultUseCase's GetHSMModule to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) GetHSMModule(ctx context.Context, input GetHSMModuleInput) (*GetHSMModuleOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.getHSMModuleInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*GetHSMModuleOutput), nil
}

// EditHSMModule implements DefaultUseCase's EditHSMModule to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) EditHSMModule(ctx context.Context, input EditHSMModuleInput) (*EditHSMModuleOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.editHSMModuleInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*EditHSMModuleOutput), nil
}

// DeleteHSMModule implements DefaultUseCase's DeleteHSMModule to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) DeleteHSMModule(ctx context.Context, input DeleteHSMModuleInput) (*DeleteHSMModuleOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.deleteHSMModuleInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*DeleteHSMModuleOutput), nil
}

func (_d *DefaultUseCaseTransactionalDecorator) createHSMModuleInternal(_ context.Context, input CreateHSMModuleInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.CreateHSMModule(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) listHSMModulesInternal(_ context.Context, generationalManagedContractListOptions ListHSMModulesInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.ListHSMModules(ctx2, generationalManagedContractListOptions)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) getHSMModuleInternal(_ context.Context, input GetHSMModuleInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.GetHSMModule(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) editHSMModuleInternal(_ context.Context, input EditHSMModuleInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.EditHSMModule(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) deleteHSMModuleInternal(_ context.Context, input DeleteHSMModuleInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.DeleteHSMModule(ctx2, input)
	}
}

var _ HSMModuleUseCase = new(DefaultUseCaseTransactionalDecorator)

// DefaultUseCaseTransactionalDecorator decorates struct DefaultUseCase wrapped with a transactional manager.
type DefaultUseCaseTransactionalDecorator struct {
	// DefaultUseCase is the usecase to be decorated.
	DefaultUseCase
	// transactionalManager defines the functionality to execute a transaction in a transactional manner.
	transactionalManager transactionalmanager.TransactionalManagerUseCase
}

// DefaultUseCaseTransactionalDecoratorOptions is the structure representing the DefaultUseCaseTransactionalDecorator dependencies.
type DefaultUseCaseTransactionalDecoratorOptions struct {
	// DefaultUseCase is the usecase to be decorated.
	DefaultUseCase *DefaultUseCase
	// TransactionalManager defines the functionality to execute a transaction in a transactional manner.
	TransactionalManager transactionalmanager.TransactionalManagerUseCase
}

// ProvideDefaultUseCaseTransactionalDecorator creates a new DefaultUseCaseTransactionalDecorator instance.
func ProvideDefaultUseCaseTransactionalDecorator(options DefaultUseCaseTransactionalDecoratorOptions) (*DefaultUseCaseTransactionalDecorator, error) {
	if options.DefaultUseCase == nil {
		errorMessage := "'DefaultAccountUseCase' is mandatory"
		return nil, errors.InvalidArgument().WithMessage(errorMessage)
	}
	if options.TransactionalManager == nil {
		errorMessage := "'TransactionalManager' is mandatory"
		return nil, errors.InvalidArgument().WithMessage(errorMessage)
	}
	return &DefaultUseCaseTransactionalDecorator{
		DefaultUseCase:       *options.DefaultUseCase,
		transactionalManager: options.TransactionalManager,
	}, nil
}
