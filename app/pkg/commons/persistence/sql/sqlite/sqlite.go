// Package sqlite defines the SQLite database functionality to operate with.
package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hyperledger-labs/signare/app/pkg/commons/persistence"
	sqlfw "github.com/hyperledger-labs/signare/app/pkg/commons/persistence/sql"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	sqlite "github.com/mattn/go-sqlite3"
)

const (
	DialectName = "sqlite"
)

func init() {
	errorMap := make(map[int]*persistence.Error)
	errorMap[19] = persistence.NewAlreadyExistsError()

	db := SQLite{
		errorMap: errorMap,
	}
	sqlfw.RegisterDialect(DialectName, &db)
}

var _ sqlfw.Dialect = new(SQLite)

const defaultMigrationsTable = "adhara_migrations"

type SQLite struct {
	db              *sqlx.DB
	conn            *sqlx.Conn
	migrationsTable string
	errorMap        map[int]*persistence.Error
}

func (s *SQLite) GetConnectionData(config any) (*sqlfw.ConnectionData, error) {
	options, ok := config.(*sqlfw.SQLiteInfo)
	if !ok {
		return nil, fmt.Errorf("invalid config to open a connection to sqlite [%v]", config)
	}

	return &sqlfw.ConnectionData{
		ConnectionString: options.ConnectionString,
		Driver:           "sqlite3",
	}, nil
}

func (s *SQLite) Connect(connectionData sqlfw.ConnectionData) (*sqlx.DB, error) {
	db, err := sqlx.Connect(connectionData.Driver, connectionData.ConnectionString)
	if err != nil {
		return nil, persistence.NewPermanentConnectionError().WithMessage(err.Error())
	}
	if db == nil {
		return nil, fmt.Errorf("cannot connect to a nil db")
	}
	db.Mapper = reflectx.NewMapperFunc("storage", func(str string) string { return str })
	s.db = db
	return db, nil
}

func (s *SQLite) OpenConnection(ctx context.Context) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *SQLite) CloseConnection(_ context.Context) error {
	return s.conn.Close()
}

func (s *SQLite) InitMigration(ctx context.Context, migrationsTablePrefix *string) error {
	if s.db == nil {
		return errors.New("you must first connect to a db")
	}
	if s.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	s.migrationsTable = defaultMigrationsTable
	if migrationsTablePrefix != nil {
		s.migrationsTable = *migrationsTablePrefix + "_" + s.migrationsTable
	}

	err := s.ensureAdharaMigrationTables(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLite) ensureAdharaMigrationTables(ctx context.Context) (err error) {
	query := `CREATE TABLE IF NOT EXISTS ` + s.migrationsTable + ` (version bigint not null primary key, dirty boolean not null, description text)`
	if _, err = s.conn.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}

func (s *SQLite) GetMigrationVersion(ctx context.Context) (*sqlfw.GetVersionOutput, error) {
	if s.db == nil {
		return nil, errors.New("you must first connect to a db")
	}
	if s.conn == nil {
		return nil, errors.New("you must first open a connection to a db")
	}

	query := `SELECT version, dirty, description FROM ` + s.migrationsTable + ` LIMIT 1`

	var version int
	var dirty bool
	var description string

	row := s.conn.QueryRowContext(ctx, query)
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

// nolint:gosec
func (s *SQLite) SetMigrationVersion(ctx context.Context, input sqlfw.SetVersionInput) error {
	if s.db == nil {
		return errors.New("you must first connect to a db")
	}
	if s.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `DELETE FROM ` + s.migrationsTable
	if _, err = tx.Exec(query); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			err = errors.Join(err, errRollback)
		}
		return err
	}

	query = `INSERT INTO ` + s.migrationsTable +
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

func (s *SQLite) RunMigration(ctx context.Context, migration string) error {
	if s.db == nil {
		return errors.New("you must first connect to a db")
	}
	if s.conn == nil {
		return errors.New("you must first open a connection to a db")
	}

	// ExecContext is not transactional, so in case of failure, some queries might have been successful and other will be pending to be executed
	_, err := s.conn.ExecContext(ctx, migration)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLite) TranslateError(_ context.Context, err error, defaultError error) error {
	var originalError sqlite.Error
	if errors.As(err, &originalError) {
		if fwError, ok := s.errorMap[int(originalError.Code)]; ok {
			return fwError
		}
	}
	return defaultError
}
