package dbtesthelper

import (
	"context"

	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/dbmigrator"
	_ "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql/init" // Used to register sql dialects
	"github.com/hyperledger-labs/signare/app/pkg/graph"
	"github.com/hyperledger-labs/signare/app/test/signaturemanagertesthelper"
)

func InitializeApp() (*graph.GraphShared, error) {
	graphConfig := graph.Config{
		BuildConfig: nil,
		Libraries: graph.LibrariesConfig{
			PersistenceFw: graph.PersistenceFwConfig{
				SQLite: &graph.SQLiteConfig{},
			},
			HSMModules: graph.HSMModules{
				SoftHSM: &graph.SoftHSMConfig{
					Library: signaturemanagertesthelper.SoftHSMLib,
				},
			},
		},
	}

	g, err := graph.New(graphConfig)
	if err != nil {
		return nil, err
	}
	g.Build()

	dbMigrator, err := dbmigrator.NewDbMigrator(dbmigrator.DbMigratorOptions{Connection: g.PersistenceFwConnection()})
	if err != nil {
		panic(err)
	}
	err = dbMigrator.MigrateFromFiles(context.Background(), dbmigrator.MigrateFromFilesInput{FS: embedded.DatabaseMigrations})
	if err != nil {
		panic(err)
	}

	return g.UseCases(), nil
}
