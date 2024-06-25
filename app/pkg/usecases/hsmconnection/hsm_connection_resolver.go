package hsmconnection

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"

	"github.com/asaskevich/govalidator"
)

// Resolver finds what HSMConnection is required depending on the constraints.
type Resolver interface {
	// ByApplication returns the HSMConnection to use by a specific application depending on its configuration. It returns an error if it fails.
	ByApplication(ctx context.Context, input ByApplicationInput) (*HSMConnection, error)
}

func (u *DefaultHSMConnectionResolver) ByApplication(ctx context.Context, input ByApplicationInput) (*HSMConnection, error) {
	_, err := govalidator.ValidateStruct(input)
	if err != nil {
		return nil, errors.InvalidArgumentFromErr(err).SetHumanReadableMessage("couldn't validate input data")
	}

	getApplicationInput := application.GetApplicationInput{
		StandardID: entities.StandardID{
			ID: input.ApplicationID,
		},
	}
	app, err := u.applicationUseCase.GetApplication(ctx, getApplicationInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	getHSMSlotInput := hsmslot.GetHSMSlotByApplicationInput{
		ApplicationID: entities.StandardID{
			ID: app.Application.ID,
		},
	}
	slot, err := u.slotUseCase.GetHSMSlotByApplication(ctx, getHSMSlotInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	getHSMModuleInput := hsmmodule.GetHSMModuleInput{
		StandardID: entities.StandardID{
			ID: slot.HSMModuleID,
		},
	}
	module, err := u.moduleUseCase.GetHSMModule(ctx, getHSMModuleInput)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}
	moduleKind, err := getModuleKind(module.Kind)
	if err != nil {
		return nil, errors.InternalFromErr(err)
	}

	return &HSMConnection{
		Slot:       slot.Slot,
		Pin:        slot.Pin,
		ChainID:    app.ChainID,
		ModuleKind: string(*moduleKind),
	}, nil
}

var _ Resolver = new(DefaultHSMConnectionResolver)

// DefaultHSMConnectionResolver implements the HSMRouter interface. It doesn't cache connections.
type DefaultHSMConnectionResolver struct {
	// moduleUseCase provides the HSM resources
	moduleUseCase hsmmodule.HSMModuleUseCase
	// slotUseCase provides the HSM slot resources
	slotUseCase hsmslot.HSMSlotUseCase
	// applicationUseCase provides the Application resources
	applicationUseCase application.ApplicationUseCase
}

// DefaultHSMConnectionResolverOptions defines options to create a new instance of DefaultHSMConnectionResolver.
type DefaultHSMConnectionResolverOptions struct {
	// ModuleUseCase provides the HSM resources
	ModuleUseCase hsmmodule.HSMModuleUseCase
	// SlotUseCase provides the HSM slot resources
	SlotUseCase hsmslot.HSMSlotUseCase
	// ApplicationUseCase provides the Application resources
	ApplicationUseCase application.ApplicationUseCase
}

// ProvideDefaultHSMConnectionResolver creates a new instance of DefaultHSMConnectionResolver using the provided options, returning an error if it fails.
func ProvideDefaultHSMConnectionResolver(options DefaultHSMConnectionResolverOptions) (*DefaultHSMConnectionResolver, error) {
	if options.ModuleUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ModuleUseCase' was not provided")
	}
	if options.SlotUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'SlotUseCase' was not provided")
	}
	if options.ApplicationUseCase == nil {
		return nil, errors.Internal().WithMessage("mandatory 'ApplicationUseCase' was not provided")
	}

	return &DefaultHSMConnectionResolver{
		moduleUseCase:      options.ModuleUseCase,
		slotUseCase:        options.SlotUseCase,
		applicationUseCase: options.ApplicationUseCase,
	}, nil
}

func getModuleKind(kind hsmmodule.ModuleKind) (*hsmconnector.ModuleKind, error) {
	var result hsmconnector.ModuleKind
	switch kind {
	case hsmmodule.SoftHSMModuleKind:
		result = hsmconnector.SoftHSMModuleKind
		return &result, nil
	default:
		return nil, errors.InvalidArgument().WithMessage("module kind '%s' not found", kind)
	}
}
