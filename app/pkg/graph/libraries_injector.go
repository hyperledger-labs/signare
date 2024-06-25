//go:build wireinject

package graph

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/wire"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"

	embedded "github.com/hyperledger-labs/signare/app"
	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
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
	if options.config.Libraries.PersistenceFw.PostgreSQL == nil || options.config.Libraries.PersistenceFw.PostgreSQL.SQLClient == nil {
		return fwOptions, nil
	}
	fwOptions.SQLClientParameters = sql.SQLClientParameters{
		MaxIdleConnections:    options.config.Libraries.PersistenceFw.PostgreSQL.SQLClient.MaxIdleConnections,
		MaxOpenConnections:    options.config.Libraries.PersistenceFw.PostgreSQL.SQLClient.MaxOpenConnections,
		MaxConnectionLifetime: options.config.Libraries.PersistenceFw.PostgreSQL.SQLClient.MaxConnectionLifetime,
	}

	return fwOptions, nil
}

func (l *librariesGraph) provideLogger(config Config) error {
	if config.Libraries.Logger == nil || config.Libraries.Logger.LogLevel == nil {
		return nil
	}
	logLevel := *config.Libraries.Logger.LogLevel
	level, ok := logger.ToLevel(logLevel)
	if !ok {
		return fmt.Errorf("level %s not defined", logLevel)
	}
	opts := logger.Options{
		Level:     level,
		LogOutput: os.Stdout,
		// Use standardized log keys to ensure interoperability and compatibility with various monitoring and tracing tools
		CtxKeysRegistry: map[entities.ContextKey]logger.LogKey{
			requestcontext.TraceParentTraceIDContextKey: "trace.id",
			requestcontext.TraceParentSpanIDContextKey:  "span.id",
			requestcontext.RPCRequestIDKey:              "rpc.request.id",
		},
	}
	logger.RegisterLogger(opts)
	return nil
}
