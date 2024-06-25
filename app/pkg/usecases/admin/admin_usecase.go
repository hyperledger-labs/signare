package admin

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/utils"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	defaultOrderDirection = entities.OrderDesc
)

const (
	signerAdminRole = "signer-admin"
)

// AdminUseCase defines the management of Admin in adminStorage.
type AdminUseCase interface {
	// CreateAdmin creates an Admin in storage and returns an error if it fails
	CreateAdmin(ctx context.Context, creation CreateAdminInput) (*CreateAdminOutput, error)
	// ListAdmins lists Admin in storage and returns an error if it fails
	ListAdmins(ctx context.Context, listOptions ListAdminsInput) (*ListAdminsOutput, error)
	// GetAdmin gets an Admin in storage and returns an error if it fails
	GetAdmin(ctx context.Context, input GetAdminInput) (*GetAdminOutput, error)
	// EditAdmin edits an Admin in storage and returns an error if it fails
	EditAdmin(ctx context.Context, update EditAdminInput) (*EditAdminOutput, error)
	// DeleteAdmin deletes an Admin in storage and returns an error if it fails
	DeleteAdmin(ctx context.Context, input DeleteAdminInput) (*DeleteAdminOutput, error)
}

func (u *DefaultUseCase) CreateAdmin(ctx context.Context, input CreateAdminInput) (*CreateAdminOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	if len(input.ID) == 0 {
		randomID := uuid.New().String()
		input.ID = randomID
	}

	now := time.Now()
	data := Admin{
		StandardResourceMeta: entities.StandardResourceMeta{
			StandardResource: entities.StandardResource{
				StandardID: input.StandardID,
				Timestamps: entities.Timestamps{
					CreationDate: now,
					LastUpdate:   now,
				},
			},
		},
		Roles:       []string{signerAdminRole},
		Description: input.Description,
	}
	data.InternalResourceID = entities.NewInternalResourceID()

	admin, err := u.adminStorage.Add(ctx, data)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err).SetHumanReadableMessage("admin [%s] already exists", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &CreateAdminOutput{
		Admin: *admin,
	}, nil
}

func (u *DefaultUseCase) ListAdmins(ctx context.Context, input ListAdminsInput) (*ListAdminsOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	filters := u.adminStorage.Filter()
	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.OrderBy == entities.OrderByLastUpdate {
		filters.OrderByLastUpdateDate(persistence.OrderDirection(direction))
	}

	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}

	collection, err := u.adminStorage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListAdminsOutput{
		AdminCollection: *collection,
	}, nil
}

func (u *DefaultUseCase) GetAdmin(ctx context.Context, input GetAdminInput) (*GetAdminOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	admin, err := u.adminStorage.Get(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("admin [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}
	return &GetAdminOutput{
		Admin: *admin,
	}, nil
}

func (u *DefaultUseCase) EditAdmin(ctx context.Context, input EditAdminInput) (*EditAdminOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	admin := Admin{
		StandardResourceMeta: input.StandardResourceMeta,
		Description:          input.Description,
		Roles:                []string{signerAdminRole},
	}
	admin.LastUpdate = time.Now()
	editedAdmin, err := u.adminStorage.Edit(ctx, admin)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("admin [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &EditAdminOutput{
		Admin: *editedAdmin,
	}, nil
}

func (u *DefaultUseCase) DeleteAdmin(ctx context.Context, input DeleteAdminInput) (*DeleteAdminOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	listAdminsInput := ListAdminsInput{}
	listAdminsOutput, err := u.ListAdmins(ctx, listAdminsInput)
	if err != nil {
		return nil, err
	}

	if len(listAdminsOutput.Items) <= 1 {
		msg := "not possible to delete admin"
		return nil, errors.PreconditionFailed().WithMessage(msg).SetHumanReadableMessage(msg)
	}

	removeAllDependenciesErr := u.removeAllDependencies(ctx, input.StandardID)
	if removeAllDependenciesErr != nil {
		return nil, removeAllDependenciesErr
	}

	admin, err := u.adminStorage.Remove(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err).SetHumanReadableMessage("admin [%s] not found", input.ID)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteAdminOutput{
		Admin: *admin,
	}, nil
}

var _ AdminUseCase = new(DefaultUseCase)

// DefaultUseCaseOptions options to create a new DefaultUseCase.
type DefaultUseCaseOptions struct {
	AdminStorage                AdminStorage
	RoleUseCase                 role.RoleUseCase
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// DefaultUseCase implementation of AdminUseCase.
type DefaultUseCase struct {
	adminStorage                AdminStorage
	roleUseCase                 role.RoleUseCase
	referentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
}

// ProvideDefaultUseCase creates a new DefaultUseCase.
func ProvideDefaultUseCase(options DefaultUseCaseOptions) (*DefaultUseCase, error) {
	if options.AdminStorage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'AdminStorage' not provided")
	}
	if options.RoleUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'RoleUseCase' not provided")
	}
	if options.ReferentialIntegrityUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ReferentialIntegrityUseCase' not provided")
	}

	return &DefaultUseCase{
		adminStorage:                options.AdminStorage,
		roleUseCase:                 options.RoleUseCase,
		referentialIntegrityUseCase: options.ReferentialIntegrityUseCase,
	}, nil
}
