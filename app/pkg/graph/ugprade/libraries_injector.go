//go:build wireinject

package upgrade

import (
	"errors"
	"os"
	"strings"

	"github.com/google/wire"

	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"
	"github.com/hyperledger-labs/signare/app/pkg/commons/time"
)

type librariesGraph struct {
	persistenceFramework  persistence.Storage
	persistenceConnection sql.Connection
}

var librariesSet = wire.NewSet(
	wire.Struct(new(librariesGraph), "*"),

	// Persistence Framework
	providePersistenceFramework,
	providePersistenceFwConnection,
	wire.Struct(new(persistenceFwConfigOptions), "*"),
	providePersistenceFwConfig,
)

func initializeLibraries(
	config Config,
) (*librariesGraph, error) {
	wire.Build(librariesSet)
	return &librariesGraph{}, nil
}

type persistenceFwConfigOptions struct {
	config Config
}

func providePersistenceFramework(connection sql.Connection, fwOptions sql.FwOptions) (persistence.Storage, error) {
	persistenceFw, err := sql.NewPersistenceFw(fwOptions)
	if err != nil {
		return nil, err
	}

	configOptions := persistence.StorageConfigOptions{
		ReadDirAndFileFS: embedded.DatabaseMappers,
		Driver:           connection.GetDialectName(),
	}

	config, err := persistence.NewStorageConfig(configOptions)
	if err != nil {
		return nil, err
	}

	err = persistenceFw.AddConfig(*config)
	if err != nil {
		return nil, err
	}

	return persistenceFw, nil
}
func providePersistenceFwConnection(options persistenceFwConfigOptions) (sql.Connection, error) {
	configs := 0
	databaseInfo := options.config.Libraries.PersistenceFw
	if databaseInfo.SQLite != nil {
		configs++
	}
	if databaseInfo.PostgreSQL != nil {
		configs++
	}
	if configs > 1 {
		return nil, errors.New("only one database can be configured at the same time")
	}

	connectionFwOptions := sql.ConnectionFwOptions{}
	if databaseInfo.PostgreSQL != nil {
		postgresConfig := sql.PostgresInfo{
			Host:     options.config.Libraries.PersistenceFw.PostgreSQL.Host,
			Port:     *options.config.Libraries.PersistenceFw.PostgreSQL.Port,
			Scheme:   *options.config.Libraries.PersistenceFw.PostgreSQL.Scheme,
			Username: options.config.Libraries.PersistenceFw.PostgreSQL.Username,
			Password: options.config.Libraries.PersistenceFw.PostgreSQL.Password,
			SSLMode:  options.config.Libraries.PersistenceFw.PostgreSQL.SSLMode,
			Database: options.config.Libraries.PersistenceFw.PostgreSQL.Database,
		}
		connectionFwOptions.Postgres = &postgresConfig
	} else if databaseInfo.SQLite != nil {
		file, err := os.CreateTemp("", strings.ReplaceAll(time.Now().String(), " ", ""))
		if err != nil {
			return nil, err
		}
		connectionFwOptions.SQLite = &sql.SQLiteInfo{
			ConnectionString: file.Name(),
		}
	} else {
		return nil, errors.New("no valid database configuration")
	}

	conn, err := sql.NewConnectionFw(connectionFwOptions)
	if err != nil {
		return sql.ConnectionFw{}, err
	}
	return conn, nil
}

func providePersistenceFwConfig(conn sql.Connection, options persistenceFwConfigOptions) (sql.FwOptions, error) {
	fwOptions := sql.FwOptions{
		Connection: conn,
	}

	return fwOptions, nil
}
