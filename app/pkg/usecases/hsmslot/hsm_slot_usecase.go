package hsmslot

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// HSMSlotUseCase defines the management of HSMSlot in storage.
type HSMSlotUseCase interface {
	// CreateHSMSlot creates a HSMSlot in storage and returns an error if it fails.
	CreateHSMSlot(ctx context.Context, input CreateHSMSlotInput) (*CreateHSMSlotOutput, error)
	// GetHSMSlot gets an HSMSlot by its ID in storage and returns an error if it fails.
	GetHSMSlot(ctx context.Context, input GetHSMSlotInput) (*GetHSMSlotOutput, error)
	// GetHSMSlotByApplication gets the HSMSlot for the specified application in storage and returns an error if it fails.
	GetHSMSlotByApplication(ctx context.Context, input GetHSMSlotByApplicationInput) (*GetHSMSlotByApplicationOutput, error)
	// EditPin edits the Pin of an HSMSlot in storage and returns an error if it fails.
	EditPin(ctx context.Context, input EditPinInput) (*EditPinOutput, error)
	// DeleteHSMSlot deletes a HSMSlot in storage and returns an error if it fails.
	DeleteHSMSlot(ctx context.Context, input DeleteHSMSlotInput) (*DeleteHSMSlotOutput, error)
	// ListHSMSlotsByApplication lists HSMSlot for a specific application in storage and returns an error if it fails.
	ListHSMSlotsByApplication(ctx context.Context, input ListHSMSlotsByApplicationInput) (*ListHSMSlotsByApplicationOutput, error)
	// ListHSMSlotsByHSMModule lists HSMSlot for a specific HSM in storage and returns an error if it fails.
	ListHSMSlotsByHSMModule(ctx context.Context, input ListHSMSlotsByHSMModuleInput) (*ListHSMSlotsByHSMModuleOutput, error)
}

func (u *DefaultUseCase) CreateHSMSlot(ctx context.Context, input CreateHSMSlotInput) (*CreateHSMSlotOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	getHSMModuleInput := hsmmodule.GetHSMModuleInput{
		StandardID: entities.StandardID{ID: input.HSMModuleID},
	}
	getHSMOutput, getHSMErr := u.hsmModuleUseCase.GetHSMModule(ctx, getHSMModuleInput)
	if getHSMErr != nil {
		if errors.IsNotFound(getHSMErr) {
			msg := fmt.Sprintf("slot '%s' in HSM '%s' is not reachable", input.Slot, input.HSMModuleID)
			return nil, errors.PreconditionFailedFromErr(getHSMErr).WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return nil, errors.InternalFromErr(getHSMErr)
	}

	resetInput := hsmconnector.ResetInput{
		ModuleKind: hsmconnector.ModuleKind(getHSMOutput.Kind),
	}
	_, resetErr := u.hsmConnector.Reset(ctx, resetInput)
	if resetErr != nil {
		return nil, resetErr
	}

	findSlotInput := hsmconnector.IsAliveInput{
		Slot:       input.Slot,
		Pin:        input.Pin,
		ModuleKind: hsmconnector.ModuleKind(getHSMOutput.Kind),
	}
	isAliveOutput, isAliveErr := u.hsmConnector.IsAlive(ctx, findSlotInput)
	if isAliveErr != nil {
		if errors.IsPreconditionFailed(isAliveErr) {
			return nil, isAliveErr
		}
		return nil, errors.InternalFromErr(isAliveErr)
	}
	if !isAliveOutput.IsAlive {
		msg := fmt.Sprintf("slot %s is not reachable in the HSM module '%s'", input.Slot, input.HSMModuleID)
		return nil, errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	hsmSlot, createSlotErr := u.createSlot(ctx, input)
	if createSlotErr != nil {
		if errors.IsAlreadyExists(createSlotErr) {
			return nil, errors.AlreadyExistsFromErr(createSlotErr)
		}
		if errors.IsPreconditionFailed(createSlotErr) {
			return nil, createSlotErr
		}
		return nil, errors.InternalFromErr(createSlotErr)
	}

	return &CreateHSMSlotOutput{
		HSMSlot: *hsmSlot,
	}, nil
}

func (u *DefaultUseCase) GetHSMSlot(ctx context.Context, input GetHSMSlotInput) (*GetHSMSlotOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	slot, err := u.hsmSlotStorage.Get(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage(fmt.Sprintf("hsm slot [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(err)
	}

	return &GetHSMSlotOutput{
		HSMSlot: *slot,
	}, nil
}

func (u *DefaultUseCase) GetHSMSlotByApplication(ctx context.Context, input GetHSMSlotByApplicationInput) (*GetHSMSlotByApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	slot, err := u.hsmSlotStorage.GetByApplication(ctx, input.ApplicationID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage(fmt.Sprintf("hsm slot not found for application [%s] not found", input.ApplicationID))
		}
		return nil, errors.InternalFromErr(err)
	}

	return &GetHSMSlotByApplicationOutput{
		HSMSlot: *slot,
	}, nil
}

func (u *DefaultUseCase) EditPin(ctx context.Context, input EditPinInput) (*EditPinOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}
	getHSMSlotInput := GetHSMSlotInput{
		StandardID: input.StandardID,
	}
	getHSMSlotOutput, getHSMSlotErr := u.GetHSMSlot(ctx, getHSMSlotInput)
	if getHSMSlotErr != nil {
		return nil, getHSMSlotErr
	}

	if getHSMSlotOutput.HSMModuleID != input.HSMModuleID {
		msg := fmt.Sprintf("slot doesn't exist in the HSM module '%s'", input.HSMModuleID)
		return nil, errors.NotFound().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	getHSMModuleInput := hsmmodule.GetHSMModuleInput{
		StandardID: entities.StandardID{ID: getHSMSlotOutput.HSMModuleID},
	}
	getHSMOutput, getHSMErr := u.hsmModuleUseCase.GetHSMModule(ctx, getHSMModuleInput)
	if getHSMErr != nil {
		if errors.IsNotFound(getHSMErr) {
			msg := fmt.Sprintf("HSM '%s' assigned to this slot does not exist", getHSMSlotOutput.HSMModuleID)
			return nil, errors.PreconditionFailedFromErr(getHSMErr).WithMessage(msg).SetHumanReadableMessage(msg)
		}
		return nil, errors.InternalFromErr(getHSMErr)
	}

	isAliveInput := hsmconnector.IsAliveInput{
		Slot:       getHSMSlotOutput.Slot,
		Pin:        input.Pin,
		ModuleKind: hsmconnector.ModuleKind(getHSMOutput.Kind),
	}

	isAliveOutput, isAliveErr := u.hsmConnector.IsAlive(ctx, isAliveInput)
	if isAliveErr != nil {
		if errors.IsPreconditionFailed(isAliveErr) {
			return nil, isAliveErr
		}
		return nil, errors.InternalFromErr(isAliveErr)
	}
	if !isAliveOutput.IsAlive {
		msg := fmt.Sprintf("slot %s is not reachable in the HSM module '%s', the new 'Pin' might be incorrect", getHSMSlotOutput.Slot, getHSMOutput.HSMModule.ID)
		return nil, errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	slot := HSMSlot{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: input.StandardID,
				Timestamps: entities.Timestamps{
					LastUpdate: time.Now(),
				},
			},
			ResourceVersion: input.ResourceVersion,
		},
		Pin: input.Pin,
	}
	editedSlot, err := u.hsmSlotStorage.EditPin(ctx, slot)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage(fmt.Sprintf("hsm slot [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(err)
	}

	return &EditPinOutput{
		HSMSlot: *editedSlot,
	}, nil
}

func (u *DefaultUseCase) DeleteHSMSlot(ctx context.Context, input DeleteHSMSlotInput) (*DeleteHSMSlotOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}
	removeAllDependenciesErr := u.removeAllDependencies(ctx, input.StandardID)
	if removeAllDependenciesErr != nil {
		return nil, removeAllDependenciesErr
	}

	removedSlot, err := u.hsmSlotStorage.Remove(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).WithMessage(fmt.Sprintf("hsm slot [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteHSMSlotOutput{
		HSMSlot: *removedSlot,
	}, nil
}

func (u *DefaultUseCase) ListHSMSlotsByApplication(ctx context.Context, input ListHSMSlotsByApplicationInput) (*ListHSMSlotsByApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.hsmSlotStorage.Filter()
	filters.FilterByApplicationID(input.ApplicationID)

	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	slotCollection, err := u.hsmSlotStorage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListHSMSlotsByApplicationOutput{
		HSMSlotCollection: *slotCollection,
	}, nil
}

func (u *DefaultUseCase) ListHSMSlotsByHSMModule(ctx context.Context, input ListHSMSlotsByHSMModuleInput) (*ListHSMSlotsByHSMModuleOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.hsmSlotStorage.Filter()
	filters.FilterByHSMModuleID(input.HSMModuleID)

	if input.ApplicationID != nil {
		filters.FilterByApplicationID(entities.StandardID{ID: *input.ApplicationID})
	}

	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	slotCollection, err := u.hsmSlotStorage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListHSMSlotsByHSMModuleOutput{
		HSMSlotCollection: *slotCollection,
	}, nil
}

func (u *DefaultUseCase) createSlot(ctx context.Context, input CreateHSMSlotInput) (*HSMSlot, error) {
	now := time.Now()
	if input.ID == nil {
		randomID := uuid.NewString()
		input.ID = &randomID
	}

	hsmSlot := HSMSlot{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: entities.StandardID{
					ID: *input.ID,
				},
				Timestamps: entities.Timestamps{
					CreationDate: now,
					LastUpdate:   now,
				},
			},
		},
		ApplicationID: input.ApplicationID,
		HSMModuleID:   input.HSMModuleID,
		Slot:          input.Slot,
		Pin:           input.Pin,
	}
	hsmSlot.InternalResourceID = entities.NewInternalResourceID()

	addApplicationDependencyErr := u.addApplicationDependency(ctx, hsmSlot)
	if addApplicationDependencyErr != nil {
		return nil, addApplicationDependencyErr
	}

	addHSMModuleDependencyErr := u.addHSMModuleDependency(ctx, hsmSlot)
	if addHSMModuleDependencyErr != nil {
		return nil, addHSMModuleDependencyErr
	}

	addedSlot, err := u.hsmSlotStorage.Add(ctx, hsmSlot)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err).SetHumanReadableMessage("hsm slot [%s] already exists", hsmSlot.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return addedSlot, nil
}

var _ HSMSlotUseCase = new(DefaultUseCase)

// DefaultUseCaseOptions configures a DefaultUseCase.
type DefaultUseCaseOptions struct {
	// HSMSlotStorage is the persistence adapter of the HSMSlot.
	HSMSlotStorage HSMSlotStorage

	// ApplicationUseCase defines how to interact with Application resources.
	ApplicationUseCase application.ApplicationUseCase
	// HSMModuleUseCase defines how to interact with HSMModule resources.
	HSMModuleUseCase hsmmodule.HSMModuleUseCase
	// HSMConnector connects with the HSM and operates with it.
	HSMConnector hsmconnector.HSMConnector
	// ReferentialIntegrityUseCase to manage dependencies between resources.
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// DefaultUseCase default management of User in configuration implementation.
type DefaultUseCase struct {
	// hsmSlotStorage is the persistence adapter of the HSMSlot.
	hsmSlotStorage HSMSlotStorage

	// applicationUseCase defines how to interact with Application resources.
	applicationUseCase application.ApplicationUseCase
	// hsmModuleUseCase defines how to interact with HSMModule resources.
	hsmModuleUseCase hsmmodule.HSMModuleUseCase
	// hsmConnector connects with the HSM and operates with it.
	hsmConnector hsmconnector.HSMConnector
	// referentialIntegrityUseCase to manage dependencies between resources.
	referentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// ProvideDefaultUseCase creates a DefaultUseCase with the given options.
func ProvideDefaultUseCase(options DefaultUseCaseOptions) (*DefaultUseCase, error) {
	if options.ApplicationUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ApplicationUseCase' was not provided")
	}
	if options.HSMConnector == nil {
		return nil, errors.Internal().WithMessage("mandatory 'HSMConnector' was not provided")
	}
	if options.HSMModuleUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'HSMModuleUseCase' was not provided")
	}
	if options.HSMSlotStorage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'HSMSlotStorage' was not provided")
	}
	if options.ReferentialIntegrityUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ReferentialIntegrityUseCase' was not provided")
	}

	return &DefaultUseCase{
		hsmConnector:                options.HSMConnector,
		hsmModuleUseCase:            options.HSMModuleUseCase,
		hsmSlotStorage:              options.HSMSlotStorage,
		applicationUseCase:          options.ApplicationUseCase,
		referentialIntegrityUseCase: options.ReferentialIntegrityUseCase,
	}, nil
}
