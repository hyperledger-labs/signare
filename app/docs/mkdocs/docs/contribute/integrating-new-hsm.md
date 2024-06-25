# Integrating with a new HSM

This document describes how to add support for a new HSM that is not yet supported, so that the signare can use it to store keys and sign data.

The target audience of the document are developers and contributors.

## Implementation

Adding support for a new HSM technology typically requires the following steps implementation wise:

1. Implementing the `DigitalSignatureManager` that defines the expected behaviour for the different HSMs.
2. Adapting the `DigitalSignatureManagerFactory` to work with the new type.
3. Extending the `HSMModule` use case to work with the new type.
4. Update the API spec with the new type.

The rest of the section will dive into those steps with a bit more detail.

Every HSM integration must implement the `DigitalSignatureManager` interface:
```
// DigitalSignatureManager defines an interface for interacting with a signature manager. It utilizes the concept of addresses to identify key pairs, where an address consists of the 20 last characters of the public key.
type DigitalSignatureManager interface {
	// GenerateKey generates a new public and private key returning the corresponding address.
	GenerateKey(ctx context.Context, input GenerateKeyInput) (*GenerateKeyOutput, error)
	// RemoveKey removes the public and private key identified by the provided address. It returns an error if it fails or if the key pair doesn't exist.
	RemoveKey(ctx context.Context, input RemoveKeyInput) (*RemoveKeyOutput, error)
	// ListKeys retrieves all stored keys as a list of addresses.
	ListKeys(ctx context.Context, input ListKeysInput) (*ListKeysOutput, error)
	// Sign signs a set of bytes with the private key identified by the provided address.
	Sign(ctx context.Context, input SignInput) (*SignOutput, error)
	// Close closes the connection and cleans up open resources.
	Close(ctx context.Context, input CloseInput) (*CloseOutput, error)
	// Open opens the connection to a digital signature manager provider.
	Open(ctx context.Context, input OpenInput) (*OpenOutput, error)
	// IsAlive checks if a given slot healthiness in a digital signature manager, returns true if it's healthy
	IsAlive(ctx context.Context, input IsAliveInput) (*IsAliveOutput, error)
}
```

Those implementations are in `app/pkg/signaturemanager`, where there is one folder for each specific HSM technology.

Additionally, all those implementations must use errors defined in the `app/pkg/signaturemanager/errors.go` file. If you need additional errors that do not exist yet, please create new ones as necessary. 

To adapt the `DigitalSignatureManagerFactory`, just register the new type along with its implementation of the `DigitalSignatureManager` interface.

Finally, extend the API specification so that resources related with the different `/admin/modules` path operations define the new type. 

## Testing

Please, refer to the documentation about [how to test](./code-standards.md#testing) for details about how to implement tests.

Every HSM must be tested following the application testing guidelines. The `HSMConnector` is the use case that uses the `DigitalSignatureManager` interface and interacts with the different HSMs. Therefore, a new test file must be created with specific tests targeting the new HSM technology, and the file must have a build tag with the name of the HSM technology. You can use the following code snippet as a template:
`hsm_connector_usecase_<hsmType>.go`
```
//go:build <hsmType>

package hsmmanager_test

func TestProvideDefaultUseCase(t *testing.T) {
    // test
}

// other tests
```

Notice that there is an `app` variable already defined in the `hsm_connector_usecase_test.go` file that contains the application graph.

## Configuration

Each HSM technology can define its own configuration attributes ([modules configuration section](../reference/configuration.md#hsm-modules-configuration){:target="_blank"}). In case that part of that configuration is dynamic and must be configured through the API, the API spec must be updated accordingly.

Additionally, `app/pkg/graph/graph_config.go` defines the configuration required by the different modules. For instance, for SoftHSM:
```
// HSMModules configures the hardware security modules.
type HSMModules struct {
	// SoftHSM configuration for SoftHSM.
	SoftHSM *SoftHSMConfig `mapstructure:"softhsm" valid:"optional"`
}
```

When adding a new module, that struct must be extended with a new attribute for the new HSM technology.

The `deployment/cmd/signare/config/config.go` file defines exposes that configuration through a static configuration file, so the equivalent data structure needs to be updated there as well. Following the previous example, for SoftHSM:
```
// HSMModules configures the hardware security modules.
type HSMModules struct {
	// SoftHSM configuration for SoftHSM.
	SoftHSM *SoftHSMConfig `mapstructure:"softhsm" valid:"optional"`
}
```
