package referentialintegrity

import (
	"context"

	"github.com/google/uuid"

	"github.com/asaskevich/govalidator"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/utils"
)

type ResourceKind string

const (
	KindAccount     ResourceKind = "account"
	KindApplication ResourceKind = "application"
	KindHSMModule   ResourceKind = "hardware_security_module"
	KindHSMSlot     ResourceKind = "hardware_security_module_slot"
	KindUser        ResourceKind = "user"
	KindAdmin       ResourceKind = "admin"
)

// ReferentialIntegrityUseCase defines how to interact with ReferentialIntegrityEntry resources.
type ReferentialIntegrityUseCase interface {
	// CreateEntry creates a ReferentialIntegrityEntry. It returns the created ReferentialIntegrityEntry or an error if it fails.
	CreateEntry(ctx context.Context, input CreateEntryInput) (*CreateEntryOutput, error)
	// GetEntry returns the requested ReferentialIntegrityEntry or an error if it fails.
	GetEntry(ctx context.Context, input GetEntryInput) (*GetEntryOutput, error)
	// ListEntries returns all the ReferentialIntegrityEntry resources or an error if it fails.
	ListEntries(ctx context.Context, input ListEntriesInput) (*ListEntriesOutput, error)
	// DeleteEntry deletes a ReferentialIntegrityEntry. It returns the deleted ReferentialIntegrityEntry or an error if it fails.
	DeleteEntry(ctx context.Context, input DeleteEntryInput) (*DeleteEntryOutput, error)
	// ListMyChildrenEntries returns all the ReferentialIntegrityEntry resources referencing a specific resource or an error if it fails.
	ListMyChildrenEntries(ctx context.Context, input ListMyChildrenEntriesInput) (*ListMyChildrenEntriesOutput, error)
	// DeleteMyEntriesIfAny deletes a ReferentialIntegrityEntry by its ID and Kind. It returns the deleted ReferentialIntegrityEntry or an error if it fails.
	DeleteMyEntriesIfAny(ctx context.Context, input DeleteMyEntriesIfAnyInput) error
	// GetEntryByResourceAndParent returns the requested ReferentialIntegrityEntry or an error if it fails.
	GetEntryByResourceAndParent(ctx context.Context, input GetEntryByResourceAndParentInput) (*GetEntryByResourceAndParentOutput, error)
}

func (u *DefaultUseCase) CreateEntry(ctx context.Context, input CreateEntryInput) (*CreateEntryOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	currentTime := time.Now()
	addInput := ReferentialIntegrityEntry{
		StandardResource: entities.StandardResource{
			StandardID: entities.StandardID{
				ID: uuid.New().String(),
			},
			Timestamps: entities.Timestamps{
				CreationDate: currentTime,
				LastUpdate:   currentTime,
			},
		},
		ResourceID:         input.ResourceID,
		ResourceKind:       input.ResourceKind,
		ParentResourceID:   input.ParentResourceID,
		ParentResourceKind: input.ParentResourceKind,
	}
	addOutput, err := u.storage.Add(ctx, addInput)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, errors.AlreadyExistsFromErr(err)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &CreateEntryOutput{
		ReferentialIntegrityEntry: *addOutput,
	}, nil
}

func (u *DefaultUseCase) GetEntry(ctx context.Context, input GetEntryInput) (*GetEntryOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	getOutput, err := u.storage.Get(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &GetEntryOutput{
		ReferentialIntegrityEntry: *getOutput,
	}, nil
}

func (u *DefaultUseCase) ListEntries(ctx context.Context, input ListEntriesInput) (*ListEntriesOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	filters := u.storage.Filter()
	direction := utils.DefaultString(input.OrderDirection, defaultOrderDirection)
	filters.OrderByCreationDate(persistence.OrderDirection(direction))
	if input.PageLimit > 0 {
		filters.Paged(input.PageLimit, input.PageOffset)
	}
	if input.Resource != nil {
		filters.FilterByResource(input.Resource.ResourceID, string(input.Resource.ResourceKind))
	}
	if input.Parent != nil {
		filters.FilterByParent(input.Parent.ResourceID, string(input.Parent.ResourceKind))
	}

	allOutput, err := u.storage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListEntriesOutput{
		ReferentialIntegrityEntryCollection: *allOutput,
	}, nil
}

func (u *DefaultUseCase) DeleteEntry(ctx context.Context, input DeleteEntryInput) (*DeleteEntryOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	removeOutput, err := u.storage.Remove(ctx, input.StandardID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFoundFromErr(err)
		}
		return nil, errors.InternalFromErr(err)
	}

	return &DeleteEntryOutput{
		ReferentialIntegrityEntry: *removeOutput,
	}, nil
}

func (u *DefaultUseCase) ListMyChildrenEntries(ctx context.Context, input ListMyChildrenEntriesInput) (*ListMyChildrenEntriesOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	filters := u.storage.Filter()
	filters.OrderByCreationDate(defaultOrderDirection)
	filters.FilterByParent(input.ParentResourceID, string(input.ParentResourceKind))

	allOutput, err := u.storage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &ListMyChildrenEntriesOutput{
		ReferentialIntegrityEntryCollection: *allOutput,
	}, nil
}

func (u *DefaultUseCase) DeleteMyEntriesIfAny(ctx context.Context, input DeleteMyEntriesIfAnyInput) error {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	_, err = u.storage.RemoveAllFromResource(ctx, input.ResourceID, string(input.ResourceKind))
	if err != nil {
		return errors.InternalFromErr(err)
	}

	return nil
}

func (u *DefaultUseCase) GetEntryByResourceAndParent(ctx context.Context, input GetEntryByResourceAndParentInput) (*GetEntryByResourceAndParentOutput, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("input data is not valid")
	}

	filters := u.storage.Filter()
	filters.FilterByParent(input.ParentResourceID, string(input.ParentResourceKind))
	filters.FilterByResource(input.ResourceID, string(input.ResourceKind))

	allOutput, err := u.storage.All(ctx, filters)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	if len(allOutput.Items) < 1 {
		return nil, errors.NotFound().WithMessage("entry with resource ID [%s], [%s], parent ID [%s] and kind [%s] not found", input.ResourceID, input.ResourceKind, input.ParentResourceID, input.ParentResourceKind)
	}

	return &GetEntryByResourceAndParentOutput{
		ReferentialIntegrityEntry: ReferentialIntegrityEntry{},
	}, nil
}

var _ ReferentialIntegrityUseCase = new(DefaultUseCase)

// DefaultUseCase default management of ReferentialIntegrityEntry in configuration implementation.
type DefaultUseCase struct {
	// storage is the persistence adapter of the references between resources.
	storage ReferentialIntegrityStorage
}

// DefaultUseCaseOptions configures a DefaultUseCase
type DefaultUseCaseOptions struct {
	// Storage is the persistence adapter of the references between resources.
	Storage ReferentialIntegrityStorage
}

// ProvideDefaultUseCase provides a DefaultUseCase with the given options
func ProvideDefaultUseCase(options DefaultUseCaseOptions) (*DefaultUseCase, error) {
	if options.Storage == nil {
		return nil, errors.Internal().WithMessage("mandatory 'Storage' was not provided")
	}
	return &DefaultUseCase{
		storage: options.Storage,
	}, nil
}
