//go:build wireinject

package graph

import (
	"github.com/google/wire"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/accountdbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/admindbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/applicationdbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/hsmdbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/hsmslotdbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/referentialintegritydbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/postgres/userdbout"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/storage/transactionaldbout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/accountdb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/admindb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/applicationdb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmmoduledb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/hsmslotdb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/referentialintegritydb"
	"github.com/hyperledger-labs/signare/app/pkg/infra/storage/userdb"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/admin"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/application"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmmodule"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/hsmslot"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/referentialintegrity"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/transactionalmanager"
	"github.com/hyperledger-labs/signare/app/pkg/usecases/user"
)

type repositoriesGraph struct {
	applicationStorage          application.ApplicationStorage
	userStorage                 user.UserStorage
	accountStorage              user.AccountStorage
	adminStorage                admin.AdminStorage
	hsmStorage                  hsmmodule.HSMModuleStorage
	hsmSlotStorage              hsmslot.HSMSlotStorage
	referentialIntegrityStorage referentialintegrity.ReferentialIntegrityStorage
	transactionalStorage        transactionalmanager.TransactionalStorage
}

var repositoriesSet = wire.NewSet(
	wire.Struct(new(repositoriesGraph), "*"),

	// Application Database Infra
	applicationdb.ProvideApplicationRepositoryInfra,
	wire.Struct(new(applicationdb.ApplicationRepositoryInfraOptions), "*"),

	// Application AdminStorage
	applicationdbout.NewRepository,
	wire.Bind(new(application.ApplicationStorage), new(*applicationdbout.Repository)),
	wire.Struct(new(applicationdbout.RepositoryOptions), "*"),

	// User Database Infra
	userdb.ProvideUserRepositoryInfra,
	wire.Struct(new(userdb.UserRepositoryInfraOptions), "*"),

	// User AdminStorage
	userdbout.NewRepository,
	wire.Bind(new(user.UserStorage), new(*userdbout.Repository)),
	wire.Struct(new(userdbout.RepositoryOptions), "*"),

	// Account Database Infra
	accountdb.ProvideAccountRepositoryInfra,
	wire.Struct(new(accountdb.AccountRepositoryInfraOptions), "*"),

	// Account AdminStorage
	accountdbout.NewRepository,
	wire.Bind(new(user.AccountStorage), new(*accountdbout.Repository)),
	wire.Struct(new(accountdbout.RepositoryOptions), "*"),

	// Admin Database Infra
	admindb.ProvideAdminRepositoryInfra,
	wire.Struct(new(admindb.AdminRepositoryInfraOptions), "*"),

	// Admin AdminStorage
	admindbout.NewRepository,
	wire.Bind(new(admin.AdminStorage), new(*admindbout.Repository)),
	wire.Struct(new(admindbout.RepositoryOptions), "*"),

	// Hardware Security Module (HSM) Database Infra
	hsmmoduledb.ProvideHardwareSecurityModuleRepositoryInfra,
	wire.Struct(new(hsmmoduledb.HardwareSecurityModuleRepositoryInfraOptions), "*"),

	// Hardware Security Module (HSM) Storage
	hsmdbout.NewRepository,
	wire.Bind(new(hsmmodule.HSMModuleStorage), new(*hsmdbout.Repository)),
	wire.Struct(new(hsmdbout.RepositoryOptions), "*"),

	// HSM Slot Database Infra
	hsmslotdb.ProvideHSMSlotRepositoryInfra,
	wire.Struct(new(hsmslotdb.HSMSlotRepositoryInfraOptions), "*"),

	// HSM Slot Storage
	hsmslotdbout.NewRepository,
	wire.Bind(new(hsmslot.HSMSlotStorage), new(*hsmslotdbout.Repository)),
	wire.Struct(new(hsmslotdbout.RepositoryOptions), "*"),

	// Referential Integrity Entry Database Infra
	referentialintegritydb.ProvideReferentialIntegrityEntryRepositoryInfra,
	wire.Struct(new(referentialintegritydb.ReferentialIntegrityEntryRepositoryInfraOptions), "*"),

	// Referential Integrity Storage
	referentialintegritydbout.NewRepository,
	wire.Bind(new(referentialintegrity.ReferentialIntegrityStorage), new(*referentialintegritydbout.Repository)),
	wire.Struct(new(referentialintegritydbout.RepositoryOptions), "*"),

	// Transactional Manager Storage
	transactionaldbout.NewTransactionalRepository,
	wire.Bind(new(transactionalmanager.TransactionalStorage), new(*transactionaldbout.TransactionalRepository)),
	wire.Struct(new(transactionaldbout.TransactionalRepositoryOptions), "*"),
)

func InitializeRepositories(
	persistenceFramework persistence.Storage,
) (*repositoriesGraph, error) {
	wire.Build(repositoriesSet)
	return &repositoriesGraph{}, nil
}
