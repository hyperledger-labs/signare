package hsmmodule

import (
	"context"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

// ModuleKind HSM module type.
type ModuleKind string

const (
	SoftHSMModuleKind ModuleKind = "SoftHSM"
)

// HSMModuleUseCase defines the management of HSMModule in storage.
type HSMModuleUseCase interface {
	// CreateHSMModule creates a CreateHSMModuleOutput in storage and returns an error if it fails.
	CreateHSMModule(ctx context.Context, creation CreateHSMModuleInput) (*CreateHSMModuleOutput, error)
	// ListHSMModules lists CreateHSMModuleOutput in storage and returns an error if it fails.
	ListHSMModules(ctx context.Context, listOptions ListHSMModulesInput) (*ListHSMModulesOutput, error)
	// GetHSMModule gets an CreateHSMModuleOutput in storage and returns an error if it fails.
	GetHSMModule(ctx context.Context, input GetHSMModuleInput) (*GetHSMModuleOutput, error)
	// EditHSMModule edits a CreateHSMModuleOutput in storage and returns an error if it fails.
	EditHSMModule(ctx context.Context, update EditHSMModuleInput) (*EditHSMModuleOutput, error)
	// DeleteHSMModule deletes a CreateHSMModuleOutput in storage and returns an error if it fails.
	DeleteHSMModule(ctx context.Context, input DeleteHSMModuleInput) (*DeleteHSMModuleOutput, error)
}

func (u *DefaultUseCase) CreateHSMModule(ctx context.Context, input CreateHSMModuleInput) (*CreateHSMModuleOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	if input.ID == nil {
		randomID := uuid.New().String()
		input.ID = &randomID
	}

	now := time.Now()
	hsmModule := HSMModule{
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
		Configuration: input.Configuration,
		Kind:          input.ModuleKind,
		Description:   input.Description,
	}
	hsmModule.InternalResourceID = entities.NewInternalResourceID()

	addedHSMModule, addErr := u.hsmModuleStorage.Add(ctx, hsmModule)
	if addErr != nil {
		if errors.IsAlreadyExists(addErr) {
			return nil, errors.AlreadyExistsFromErr(addErr).SetHumanReadableMessage("hsm module [%s] already exists", *input.ID)
		}
		return nil, errors.InternalFromErr(addErr)
	}

	return &CreateHSMModuleOutput{
		HSMModule: *addedHSMModule,
	}, nil
}

func (u *DefaultUseCase) ListHSMModules(ctx context.Context, input ListHSMModulesInput) (*ListHSMModulesOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.hsmModuleStorage.Filter()
	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	hsmModulesCollection, err := u.hsmModuleStorage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	return &ListHSMModulesOutput{
		HSMModulesCollection: *hsmModulesCollection,
	}, nil
}

func (u *DefaultUseCase) GetHSMModule(ctx context.Context, input GetHSMModuleInput) (*GetHSMModuleOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	hsmModule, getHsmModuleErr := u.hsmModuleStorage.Get(ctx, input.StandardID)
	if getHsmModuleErr != nil {
		if errors.IsNotFound(getHsmModuleErr) {
			return nil, errors.NotFoundFromErr(getHsmModuleErr).WithMessage(fmt.Sprintf("hsm module [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(getHsmModuleErr)
	}

	return &GetHSMModuleOutput{
		HSMModule: *hsmModule,
	}, nil
}

func (u *DefaultUseCase) EditHSMModule(ctx context.Context, input EditHSMModuleInput) (*EditHSMModuleOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	input.HSMModule.LastUpdate = time.Now()
	hsmModule, editHSMModuleErr := u.hsmModuleStorage.Edit(ctx, input.HSMModule)
	if editHSMModuleErr != nil {
		if errors.IsNotFound(editHSMModuleErr) {
			return nil, errors.NotFoundFromErr(editHSMModuleErr).WithMessage(fmt.Sprintf("hsm module [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(editHSMModuleErr)
	}

	return &EditHSMModuleOutput{
		HSMModule: *hsmModule,
	}, nil
}

func (u *DefaultUseCase) DeleteHSMModule(ctx context.Context, input DeleteHSMModuleInput) (*DeleteHSMModuleOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	removeAllDependenciesErr := u.removeAllDependencies(ctx, input.StandardID)
	if removeAllDependenciesErr != nil {
		return nil, removeAllDependenciesErr
	}

	hsmModule, removeHSMModuleErr := u.hsmModuleStorage.Remove(ctx, input.StandardID)
	if removeHSMModuleErr != nil {
		if errors.IsNotFound(removeHSMModuleErr) {
			return nil, errors.NotFoundFromErr(removeHSMModuleErr).WithMessage(fmt.Sprintf("hsm module [%s] not found", input.ID))
		}
		return nil, errors.InternalFromErr(removeHSMModuleErr)
	}

	return &DeleteHSMModuleOutput{
		HSMModule: *hsmModule,
	}, nil
}

var _ HSMModuleUseCase = new(DefaultUseCase)

// DefaultUseCaseOptions options to create a DefaultUseCase.
type DefaultUseCaseOptions struct {
	// HSMModuleStorage is the persistence adapter of the HSMModule.
	HSMModuleStorage HSMModuleStorage
	// ReferentialIntegrityUseCase to manage dependencies between resources.
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// DefaultUseCase implements the HSMModuleUseCase interface.
type DefaultUseCase struct {
	// hsmModuleStorage is the persistence adapter of the HSMModule.
	hsmModuleStorage HSMModuleStorage
	// referentialIntegrityUseCase to manage dependencies between resources.
	referentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// ProvideDefaultHSMModuleUseCase creates a new DefaultUseCase instance.
func ProvideDefaultHSMModuleUseCase(options DefaultUseCaseOptions) (*DefaultUseCase, error) {
	if options.HSMModuleStorage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'HSMModuleUseCase' not provided")
	}
	if options.ReferentialIntegrityUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ReferentialIntegrityUseCase' not provided")
	}

	return &DefaultUseCase{
		hsmModuleStorage:            options.HSMModuleStorage,
		referentialIntegrityUseCase: options.ReferentialIntegrityUseCase,
	}, nil
}
