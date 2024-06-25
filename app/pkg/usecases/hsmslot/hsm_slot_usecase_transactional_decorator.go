package hsmslot

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"
)

// CreateHSMSlot implements DefaultUseCase's CreateHSMSlot to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) CreateHSMSlot(ctx context.Context, input CreateHSMSlotInput) (*CreateHSMSlotOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.createHSMSlotInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*CreateHSMSlotOutput), nil
}

// GetHSMSlot implements DefaultUseCase's GetHSMSlot to be a transactional operation
func (_d *DefaultUseCaseTransactionalDecorator) GetHSMSlot(ctx context.Context, input GetHSMSlotInput) (*GetHSMSlotOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.getHSMSlotInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*GetHSMSlotOutput), nil
}

// EditPin implements DefaultUseCase's EditPin to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) EditPin(ctx context.Context, input EditPinInput) (*EditPinOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.editPinInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*EditPinOutput), nil
}

// DeleteHSMSlot implements DefaultUseCase's DeleteHSMSlot to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) DeleteHSMSlot(ctx context.Context, input DeleteHSMSlotInput) (*DeleteHSMSlotOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.deleteHSMSlot(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*DeleteHSMSlotOutput), nil
}

// ListHSMSlotsByApplication implements DefaultUseCase's ListHSMSlotsByApplication to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) ListHSMSlotsByApplication(ctx context.Context, input ListHSMSlotsByApplicationInput) (*ListHSMSlotsByApplicationOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.listHSMSlotsByApplicationInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*ListHSMSlotsByApplicationOutput), nil
}

// ListHSMSlotsByHSMModule implements DefaultUseCase's ListHSMSlotsByHSMModule to be a transactional operation.
func (_d *DefaultUseCaseTransactionalDecorator) ListHSMSlotsByHSMModule(ctx context.Context, input ListHSMSlotsByHSMModuleInput) (*ListHSMSlotsByHSMModuleOutput, error) {
	returnValue, failure := _d.transactionalManager.ExecuteInTransaction(ctx, _d.listHSMSlotsByHSMModuleInternal(ctx, input))
	if failure != nil {
		return nil, failure
	}

	return returnValue.(*ListHSMSlotsByHSMModuleOutput), nil
}

func (_d *DefaultUseCaseTransactionalDecorator) createHSMSlotInternal(_ context.Context, input CreateHSMSlotInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.CreateHSMSlot(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) getHSMSlotInternal(_ context.Context, generationalManagedContractListOptions GetHSMSlotInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.GetHSMSlot(ctx2, generationalManagedContractListOptions)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) editPinInternal(_ context.Context, input EditPinInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.EditPin(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) deleteHSMSlot(_ context.Context, input DeleteHSMSlotInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.DeleteHSMSlot(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) listHSMSlotsByApplicationInternal(_ context.Context, input ListHSMSlotsByApplicationInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.ListHSMSlotsByApplication(ctx2, input)
	}
}

func (_d *DefaultUseCaseTransactionalDecorator) listHSMSlotsByHSMModuleInternal(_ context.Context, input ListHSMSlotsByHSMModuleInput) func(context.Context) (interface{}, error) {
	return func(ctx2 context.Context) (interface{}, error) {
		return _d.DefaultUseCase.ListHSMSlotsByHSMModule(ctx2, input)
	}
}

var _ HSMSlotUseCase = new(DefaultUseCaseTransactionalDecorator)

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

// ProvideDefaultUseCaseTransactionalDecorator creates a new DefaultUseCaseTransactionalDecorator.
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
