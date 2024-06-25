package hsmconnector

import (
	"context"
	"fmt"
	"os"

	"github.com/miekg/pkcs11"

	signererrors "github.com/hyperledger-labs/signare/app/pkg/internal/errors"
	"github.com/hyperledger-labs/signare/app/pkg/signaturemanager"
	"github.com/hyperledger-labs/signare/app/pkg/signaturemanager/pkcs11hsm"
)

// DigitalSignatureManagerFactory defines the factory to create DigitalSignatureManager connections.
type DigitalSignatureManagerFactory interface {
	// Create returns the connection to a DigitalSignatureManager.
	Create(ctx context.Context, input CreateInput) (signaturemanager.DigitalSignatureManager, error)
	// Close closes open resources to the digital signature manager.
	Close(ctx context.Context, input CloseInput) (*CloseOutput, error)
	// Reset the snapshot of a given module kind to include slots created after the initialization.
	Reset(_ context.Context, kind ModuleKind) error
}

func (u *DefaultDigitalSignatureManagerFactory) Reset(ctx context.Context, kind ModuleKind) error {
	digitalSignatureManager, ok := u.digitalSignatureManagerMap[kind]
	if !ok {
		return signererrors.InvalidArgument().WithMessage("error during reset of the digital signature manager: the HSM type '%s' is not supported", kind)
	}
	_, closeErr := digitalSignatureManager.Close(ctx, signaturemanager.CloseInput{})
	if closeErr != nil {
		return signererrors.Internal().WithMessage("error closing digital signature manager connection '%s'. Error: %v", kind, closeErr)
	}
	_, openErr := digitalSignatureManager.Open(ctx, signaturemanager.OpenInput{})
	if openErr != nil {
		return signererrors.Internal().WithMessage("error opening digital signature manager connection '%s'. Error: %v", kind, openErr)
	}
	u.digitalSignatureManagerMap[kind] = digitalSignatureManager

	return nil
}

func (u *DefaultDigitalSignatureManagerFactory) Create(ctx context.Context, input CreateInput) (signaturemanager.DigitalSignatureManager, error) {
	if input.ModuleKind != SoftHSMModuleKind {
		errMsg := fmt.Sprintf("the provided module kind '%s' is not supported", input.ModuleKind)
		return nil, signererrors.InvalidArgument().SetHumanReadableMessage(errMsg).WithMessage(errMsg)
	}

	digitalSignatureManager, ok := u.digitalSignatureManagerMap[input.ModuleKind]
	if !ok {
		return nil, signererrors.InvalidArgument().WithMessage("the provided module kind '%s' is not supported", input.ModuleKind)
	}

	_, openErr := digitalSignatureManager.Open(ctx, signaturemanager.OpenInput{})
	if openErr != nil && !signaturemanager.IsAlreadyInitializedErr(openErr) {
		return nil, signererrors.Internal().WithMessage("failed to open digital signature manager '%s'. Error: %v", input.ModuleKind, openErr)
	}

	return digitalSignatureManager, nil
}

func (u *DefaultDigitalSignatureManagerFactory) Close(ctx context.Context, _ CloseInput) (*CloseOutput, error) {
	for key, digitalSignatureManager := range u.digitalSignatureManagerMap {
		_, err := digitalSignatureManager.Close(ctx, signaturemanager.CloseInput{})
		if err != nil {
			return nil, signererrors.InternalFromErr(err).WithMessage("error closing digital signature manager: '%s'. Error: %v", key, err)
		}
	}
	return &CloseOutput{}, nil
}

var _ DigitalSignatureManagerFactory = new(DefaultDigitalSignatureManagerFactory)

// DefaultDigitalSignatureManagerFactory implements DigitalSignatureManagerFactory to create PKCS11 digital signature
// manager compatible instances.
// It Initializes the pkcs11 library at creation time so that there is one pkcs11.Ctx per digital signature manager supported type.
type DefaultDigitalSignatureManagerFactory struct {
	digitalSignatureManagerMap map[ModuleKind]signaturemanager.DigitalSignatureManager
}

// DefaultDigitalSignatureManagerFactoryOptions options to create a new DigitalSignatureManagerFactory instance.
type DefaultDigitalSignatureManagerFactoryOptions struct {
	// SoftHSMLibrary path to the library to connect to a PKCS11 compatible HSM.
	SoftHSMLibrary *PKCS11Library
}

// ProvideDefaultDigitalSignatureManagerFactory creates a new DigitalSignatureManagerFactory with the given options.
func ProvideDefaultDigitalSignatureManagerFactory(options DefaultDigitalSignatureManagerFactoryOptions) (*DefaultDigitalSignatureManagerFactory, error) {
	digitalSignatureManagerMap := make(map[ModuleKind]signaturemanager.DigitalSignatureManager)
	if options.SoftHSMLibrary != nil {
		_, err := os.Stat(string(*options.SoftHSMLibrary))
		if os.IsNotExist(err) {
			return nil, signererrors.InvalidArgument().WithMessage("SoftHSM library path does not exist")
		}
		pkcs11Context := pkcs11.New(string(*options.SoftHSMLibrary))
		if pkcs11Context == nil {
			return nil, signererrors.Internal().WithMessage("error instantiating the PKCS11 interface for '%s'", SoftHSMModuleKind)
		}
		errInitialize := pkcs11Context.Initialize()
		if errInitialize != nil {
			return nil, signererrors.Internal().WithMessage("error calling the PKCS11 interface initialize function for '%s'. Error: %v", SoftHSMModuleKind, errInitialize)
		}
		pkcs11HSMSignatureManagerOptions := pkcs11hsm.PKCS11HSMSignatureManagerOptions{
			PkcsContext: pkcs11Context,
		}
		signatureManager, err := pkcs11hsm.ProvidePKCS11HSMSignatureManager(pkcs11HSMSignatureManagerOptions)
		if err != nil {
			return nil, signererrors.InternalFromErr(err)
		}
		digitalSignatureManagerMap[SoftHSMModuleKind] = signatureManager
	}

	if len(digitalSignatureManagerMap) == 0 {
		return nil, signererrors.InvalidArgument().WithMessage("no HSM libraries were provided. At least one is required")
	}

	return &DefaultDigitalSignatureManagerFactory{
		digitalSignatureManagerMap: digitalSignatureManagerMap,
	}, nil
}
