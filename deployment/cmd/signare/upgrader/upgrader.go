// Package upgrader defines the upgrade database command.
package upgrader

import (
	"context"
	"embed"
	"fmt"

	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/dbmigrator"
	_ "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql/init" // Used to register sql dialects
	upgrade "github.com/hyperledger-labs/signare/app/pkg/graph/ugprade"
	"github.com/hyperledger-labs/signare/deployment/cmd/signare/config"
	"github.com/hyperledger-labs/signare/deployment/cmd/signare/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "upgrade",
		Long: "upgrades to the given version",
		RunE: executeUpgrade,
	}
	return cmd
}

func executeUpgrade(_ *cobra.Command, _ []string) error {
	ctx := context.Background()
	configFilePath := viper.GetString(flags.SignareConfigPathFlag)
	staticConfig, err := config.GetStaticConfiguration(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("error reading static configuration: [%v]", err))
	}

	appConfig := toGraphConfiguration(staticConfig)
	appGraph, err := upgrade.New(appConfig)
	if err != nil {
		panic(fmt.Sprintf("error initializing appGraph: [%v]", err))
	}

	appGraph.Build()
	upgradeDatabase(ctx, embedded.DatabaseMigrations, *appGraph)

	return nil
}

func toGraphConfiguration(staticConfig *config.StaticConfiguration) upgrade.Config {
	graphConfig := upgrade.Config{
		Libraries: upgrade.LibrariesConfig{
			PersistenceFw: upgrade.PersistenceFwConfig{
				PostgreSQL: &upgrade.PostgresSQLConfig{
					Host:     staticConfig.DatabaseInfo.PostgreSQL.Host,
					Port:     &staticConfig.DatabaseInfo.PostgreSQL.Port,
					Scheme:   &staticConfig.DatabaseInfo.PostgreSQL.Scheme,
					Username: staticConfig.DatabaseInfo.PostgreSQL.Username,
					Password: staticConfig.DatabaseInfo.PostgreSQL.Password,
					SSLMode:  staticConfig.DatabaseInfo.PostgreSQL.SSLMode,
					Database: staticConfig.DatabaseInfo.PostgreSQL.Database,
				},
			},
		},
	}

	return graphConfig
}

func upgradeDatabase(ctx context.Context, fs embed.FS, upgradeGraph upgrade.UpgradeGraph) {
	dbMigrator, err := dbmigrator.NewDbMigrator(dbmigrator.DbMigratorOptions{Connection: upgradeGraph.PersistenceFwConnection()})
	if err != nil {
		panic(err)
	}
	err = dbMigrator.MigrateFromFiles(ctx, dbmigrator.MigrateFromFilesInput{FS: fs})
	if err != nil {
		panic(err)
	}
}
