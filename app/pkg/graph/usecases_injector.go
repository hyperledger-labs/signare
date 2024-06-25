//go:build wireinject

package graph

import (
	"github.com/google/wire"

	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnector"

	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"

	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/infile/roleinfile"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/authorization/role"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmconnection"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

type useCasesGraph struct {
	ApplicationUseCase          application.ApplicationUseCase
	UserUseCase                 user.UserUseCase
	AccountUseCase              user.AccountUseCase
	AdminUseCase                admin.AdminUseCase
	HSMModuleUseCase            hsmmodule.HSMModuleUseCase
	HSMSlotUseCase              hsmslot.HSMSlotUseCase
	HSMConnector                hsmconnector.HSMConnector
	RoleUseCase                 role.RoleUseCase
	HSMConnectionResolver       hsmconnection.Resolver
	ReferentialIntegrityUseCase referentialintegrity.ReferentialIntegrityUseCase
	TransactionalManagerUseCase transactionalmanager.TransactionalManagerUseCase

	DigitalSignatureManagerFactory hsmconnector.DigitalSignatureManagerFactory
}

var useCasesSet = wire.NewSet(
	wire.Struct(new(useCasesGraph), "*"),

	// Transactional Manager Use Case
	transactionalmanager.ProvideTransactionalManager,
	wire.Bind(new(transactionalmanager.TransactionalManagerUseCase), new(*transactionalmanager.TransactionalManager)),
	wire.Struct(new(transactionalmanager.TransactionalManagerOptions), "*"),

	// Referential Integrity Use Case
	referentialintegrity.ProvideDefaultUseCase,
	wire.Bind(new(referentialintegrity.ReferentialIntegrityUseCase), new(*referentialintegrity.DefaultUseCase)),
	wire.Struct(new(referentialintegrity.DefaultUseCaseOptions), "*"),

	// Application Use Case
	application.ProvideDefaultUseCase,
	wire.Bind(new(application.ApplicationUseCase), new(*application.DefaultUseCase)),
	wire.Struct(new(application.DefaultUseCaseOptions), "*"),

	// User Use Case
	user.ProvideDefaultUseCase,
	wire.Bind(new(user.UserUseCase), new(*user.DefaultUserUseCase)),
	wire.Struct(new(user.DefaultUserUseCaseOptions), "*"),

	// Account Use Case [Transactional]
	user.ProvideDefaultUseCaseTransactionalDecorator,
	wire.Bind(new(user.AccountUseCase), new(*user.DefaultUserUseCase)),
	wire.Struct(new(user.DefaultUseCaseTransactionalDecoratorOptions), "*"),

	// Admin Use Case
	admin.ProvideDefaultUseCase,
	wire.Bind(new(admin.AdminUseCase), new(*admin.DefaultUseCase)),
	wire.Struct(new(admin.DefaultUseCaseOptions), "*"),

	// HSM Module Use Case [Transactional]
	hsmmodule.ProvideDefaultUseCaseTransactionalDecorator,
	wire.Bind(new(hsmmodule.HSMModuleUseCase), new(*hsmmodule.DefaultUseCaseTransactionalDecorator)),
	wire.Struct(new(hsmmodule.DefaultUseCaseTransactionalDecoratorOptions), "*"),
	hsmmodule.ProvideDefaultHSMModuleUseCase,
	wire.Struct(new(hsmmodule.DefaultUseCaseOptions), "*"),

	// HSM Slot Use Case [Transactional]
	hsmslot.ProvideDefaultUseCaseTransactionalDecorator,
	wire.Bind(new(hsmslot.HSMSlotUseCase), new(*hsmslot.DefaultUseCaseTransactionalDecorator)),
	wire.Struct(new(hsmslot.DefaultUseCaseTransactionalDecoratorOptions), "*"),
	hsmslot.ProvideDefaultUseCase,
	wire.Struct(new(hsmslot.DefaultUseCaseOptions), "*"),

	// HMS Connector Use Case
	hsmconnector.ProvideDefaultHSMConnector,
	wire.Bind(new(hsmconnector.HSMConnector), new(*hsmconnector.DefaultUseCase)),
	wire.Struct(new(hsmconnector.DefaultUseCaseOptions), "*"),

	// Role Use Case
	provideDefaultRoleStorageInFile,
	role.ProvideDefaultRoleUseCase,
	wire.Bind(new(role.RoleUseCase), new(*role.DefaultRoleUseCase)),
	wire.Struct(new(role.DefaultRoleUseCaseOptions), "*"),

	// Digital Signature Manager DigitalSignatureManagerFactory
	provideSoftHSMConfiguration,
	hsmconnector.ProvideDefaultDigitalSignatureManagerFactory,
	wire.Bind(new(hsmconnector.DigitalSignatureManagerFactory), new(*hsmconnector.DefaultDigitalSignatureManagerFactory)),
	wire.Struct(new(hsmconnector.DefaultDigitalSignatureManagerFactoryOptions), "*"),

	// HSM Connection
	hsmconnection.ProvideDefaultHSMConnectionResolver,
	wire.Bind(new(hsmconnection.Resolver), new(*hsmconnection.DefaultHSMConnectionResolver)),
	wire.Struct(new(hsmconnection.DefaultHSMConnectionResolverOptions), "*"),
)

func initializeUseCases(
	repositories *repositoriesGraph,
	metricRecorder metricrecorder.MetricRecorder,
	config Config,
) (*useCasesGraph, error) {
	wire.Build(useCasesSet,
		wire.FieldsOf(new(*repositoriesGraph),
			"applicationStorage",
			"userStorage",
			"accountStorage",
			"adminStorage",
			"hsmStorage",
			"hsmSlotStorage",
			"referentialIntegrityStorage",
			"transactionalStorage",
		),
	)
	return &useCasesGraph{}, nil
}

func provideSoftHSMConfiguration(config Config) *hsmconnector.PKCS11Library {
	if config.Libraries.HSMModules.SoftHSM != nil {
		softHSMConfig := hsmconnector.PKCS11Library(config.Libraries.HSMModules.SoftHSM.Library)
		return &softHSMConfig
	}
	return nil
}

func provideDefaultRoleStorageInFile() role.RoleStorage {
	defaultRoleStorageInFileOptions := roleinfile.DefaultRoleStorageInFileOptions{
		FileSystem: embedded.RBACFiles,
		BasePath:   "include/rbac",
	}
	defaultRoleStorageInFile, err := roleinfile.ProvideDefaultRoleStorageInFile(defaultRoleStorageInFileOptions)
	checkError(err)

	return defaultRoleStorageInFile
}
