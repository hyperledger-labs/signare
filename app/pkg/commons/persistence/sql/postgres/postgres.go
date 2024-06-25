// Package postgres defines the PostgresSQL functionalities.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strconv"
	"syscall"

	"github.com/jmoiron/sqlx/reflectx"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // register postgresql driver
	"github.com/jmoiron/sqlx"

	sqlfw "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"
)

const (
	DialectName = "postgres"
)

func init() {
	errorMap := make(map[string]*persistence.Error)
	errorMap["23505"] = persistence.NewAlreadyExistsError()
	errorMap["3D000"] = persistence.NewTransientConnectionError().WithMessage("couldn't find the given table in the PostgreSQL server")
	errorMap["28P01"] = persistence.NewPermanentConnectionError().WithMessage("forbidden access")

	db := Postgres{
		errorMap: errorMap,
	}
	sqlfw.RegisterDialect(DialectName, &db)
}

const defaultMigrationsTable = "adhara_migrations"

type Postgres struct {
	db              *sqlx.DB
	conn            *sqlx.Conn
	migrationsTable string
	errorMap        map[string]*persistence.Error
}

func (p *Postgres) GetConnectionData(config interface{}) (*sqlfw.ConnectionData, error) {
	options, ok := config.(*sqlfw.PostgresInfo)
	if !ok {
		return nil, fmt.Errorf("invalid config to open a connection to postgres [%v]", config)
	}

	port := strconv.Itoa(options.Port)

	if options.Username == "" {
		return nil, errors.New("variable 'Username' cannot be empty")
	}
	if options.Host == "" {
		return nil, errors.New("variable 'Host' cannot be empty")
	}
	if port == "" {
		return nil, errors.New("variable 'Port' cannot be empty")
	}
	if options.Database == "" {
		return nil, errors.New("variable 'Database' cannot be empty")
	}
	if options.SSLMode == "" {
		return nil, errors.New("variable 'SSLMode' cannot be empty")
	}

	connectionString := fmt.Sprintf("application_name=SIGNARE host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", options.Host, port, options.Username, options.Password, options.Database, options.SSLMode)

	return &sqlfw.ConnectionData{
		ConnectionString: connectionString,
		Driver:           "pgx",
	}, nil
}

func (p *Postgres) Connect(connectionData sqlfw.ConnectionData) (*sqlx.DB, error) {
	db, err := sqlx.Connect(connectionData.Driver, connectionData.ConnectionString)
	if err != nil {
		return nil, p.translateConnectionError(context.Background(), err, persistence.NewPermanentConnectionError().WithMessage(err.Error()))
	}
	if db == nil {
		return nil, errors.New("cannot connect to a nil db")
	}
	db.Mapper = reflectx.NewMapperFunc("storage", func(str string) string { return str })
	p.db = db
	return db, nil
}

func (p *Postgres) OpenConnection(ctx context.Context) error {
	conn, err := p.db.Connx(ctx)
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

func (p *Postgres) CloseConnection(_ context.Context) error {
	return p.conn.Close()
}

func (p *Postgres) InitMigration(ctx context.Context, migrationsTablePrefix *string) error {
	if p.db == nil {
		return errors.New("you must first connect to a db")
	}
	if p.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	p.migrationsTable = defaultMigrationsTable
	if migrationsTablePrefix != nil {
		p.migrationsTable = *migrationsTablePrefix + "_" + p.migrationsTable
	}

	err := p.ensureAdharaMigrationTables(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) SetMigrationVersion(ctx context.Context, input sqlfw.SetVersionInput) error {
	if p.db == nil {
		return errors.New("you must first connect to a db")
	}
	if p.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	tx, err := p.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `TRUNCATE ` + p.migrationsTable
	if _, err = tx.Exec(query); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = errors.Join(err, errRollback)
		}
		return err
	}

	query = `INSERT INTO ` + p.migrationsTable +
		` (version, dirty, description) VALUES ($1, $2, $3)`
	if _, err = tx.Exec(query, input.Version, input.Dirty, input.Description); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = errors.Join(err, errRollback)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetMigrationVersion(ctx context.Context) (*sqlfw.GetVersionOutput, error) {
	if p.db == nil {
		return nil, errors.New("you must first connect to a db")
	}
	if p.conn == nil {
		return nil, errors.New("you must first open a connection to a db")
	}

	query := `SELECT version, dirty, description FROM ` + p.migrationsTable + ` LIMIT 1`

	var version int
	var dirty bool
	var description string

	row := p.conn.QueryRowContext(ctx, query)
	err := row.Scan(&version, &dirty, &description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &sqlfw.GetVersionOutput{
				MigrationVersion: sqlfw.MigrationVersion{
					Version: 0,
					Dirty:   false,
				},
			}, nil
		}
		return nil, err
	}

	return &sqlfw.GetVersionOutput{
		MigrationVersion: sqlfw.MigrationVersion{
			Version:     version,
			Dirty:       dirty,
			Description: description,
		},
	}, nil
}

func (p *Postgres) RunMigration(ctx context.Context, migration string) error {
	if p.db == nil {
		return errors.New("you must first connect to a db")
	}
	if p.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	// ExecContext is not transactional, so in case of failure, some queries might have been successful and other will be pending to be executed
	_, err := p.conn.ExecContext(ctx, migration)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) ensureAdharaMigrationTables(ctx context.Context) (err error) {
	query := `CREATE TABLE IF NOT EXISTS ` + p.migrationsTable + ` (version bigint not null primary key, dirty boolean not null, description text)`
	if _, err = p.conn.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) TranslateError(_ context.Context, err error, defaultError error) error {
	var originalError *pgconn.PgError
	if errors.As(err, &originalError) {
		if fwError, ok := p.errorMap[originalError.Code]; ok {
			return fwError
		}
	}
	return defaultError
}

func (p *Postgres) translateConnectionError(ctx context.Context, err error, defaultError error) error {
	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr != nil {
		return p.translateConnectionError(ctx, unwrappedErr, defaultError)
	}
	var sysCallErr syscall.Errno
	okSysCallErr := errors.As(err, &sysCallErr)
	if okSysCallErr {
		if sysCallErr.Temporary() {
			return persistence.NewTransientConnectionError().WithMessage(err.Error())
		}
		if errors.Is(sysCallErr, syscall.ECONNREFUSED) {
			return persistence.NewTransientConnectionError().WithMessage("PostgreSQL server not found in given port")
		}

		return persistence.NewPermanentConnectionError().WithMessage(err.Error())
	}

	var netDNSErrorErr *net.DNSError
	okDNSErr := errors.As(err, &netDNSErrorErr)
	if okDNSErr {
		if netDNSErrorErr.Temporary() || netDNSErrorErr.Timeout() {
			return persistence.NewTransientConnectionError().WithMessage(err.Error())
		}

		return persistence.NewPermanentConnectionError().WithMessage(err.Error())
	}

	return p.TranslateError(ctx, err, defaultError)
}

var _ sqlfw.Dialect = (*Postgres)(nil)
