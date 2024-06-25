package application

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/asaskevich/govalidator"
)

// ApplicationUseCase defines how to interact with Application resources.
type ApplicationUseCase interface {
	// CreateApplication creates an Application in configuration and returns an error if it fails
	CreateApplication(ctx context.Context, input CreateApplicationInput) (*CreateApplicationOutput, error)
	// ListApplications lists Applications in configuration and returns an error if it fails
	ListApplications(ctx context.Context, input ListApplicationsInput) (*ListApplicationsOutput, error)
	// GetApplication gets an Application in configuration and returns an error if it fails
	GetApplication(ctx context.Context, input GetApplicationInput) (*GetApplicationOutput, error)
	// EditApplication edits an Application in configuration and returns an error if it fails
	EditApplication(ctx context.Context, input EditApplicationInput) (*EditApplicationOutput, error)
	// DeleteApplication deletes an Application in configuration and returns an error if it fails
	DeleteApplication(ctx context.Context, input DeleteApplicationInput) (*DeleteApplicationOutput, error)
}

var _ ApplicationUseCase = (*DefaultUseCase)(nil)

func (u *DefaultUseCase) CreateApplication(ctx context.Context, input CreateApplicationInput) (*CreateApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	appToCreate := createToApplication(input)
	appToCreate.InternalResourceID = entities.NewInternalResourceID()

	application, err := u.storage.Add(ctx, appToCreate)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err).SetHumanReadableMessage("application [%s] already exists", appToCreate.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &CreateApplicationOutput{
		Application: *application,
	}, nil
}

func (u *DefaultUseCase) ListApplications(ctx context.Context, input ListApplicationsInput) (*ListApplicationsOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.storage.Filter()
	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	applicationCollection, err := u.storage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	return &ListApplicationsOutput{
		ApplicationCollection: *applicationCollection,
	}, nil
}

func (u *DefaultUseCase) GetApplication(ctx context.Context, input GetApplicationInput) (*GetApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	getInput := entities.StandardID{
		ID: input.ID,
	}
	application, err := u.storage.Get(ctx, getInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("application [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}
	return &GetApplicationOutput{
		Application: *application,
	}, nil
}

func (u *DefaultUseCase) EditApplication(ctx context.Context, input EditApplicationInput) (*EditApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	application := mapEditedValues(input)
	editedApplication, err := u.storage.Edit(ctx, application)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("application [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &EditApplicationOutput{
		Application: *editedApplication,
	}, nil
}

func (u *DefaultUseCase) DeleteApplication(ctx context.Context, input DeleteApplicationInput) (*DeleteApplicationOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	removeAllDependenciesErr := u.removeAllDependencies(ctx, input.StandardID)
	if removeAllDependenciesErr != nil {
		return nil, removeAllDependenciesErr
	}

	removeInput := entities.StandardID{
		ID: input.ID,
	}
	application, err := u.storage.Remove(ctx, removeInput)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("application [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}
	return &DeleteApplicationOutput{
		Application: *application,
	}, nil
}

var _ ApplicationUseCase = new(DefaultUseCase)

// DefaultUseCase default management of Application in configuration implementation
type DefaultUseCase struct {
	storage                     ApplicationStorage
	referentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// DefaultUseCaseOptions configures a DefaultUseCase
type DefaultUseCaseOptions struct {
	Storage                     ApplicationStorage
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// ProvideDefaultUseCase provides a DefaultUseCase with the given options
func ProvideDefaultUseCase(options DefaultUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'Storage' not provided")
	}
	if options.ReferentialIntegrityUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ReferentialIntegrityUseCase' not provided")
	}

	return &DefaultUseCase{
		storage:                     options.Storage,
		referentialIntegrityUseCase: options.ReferentialIntegrityUseCase,
	}, nil
}
