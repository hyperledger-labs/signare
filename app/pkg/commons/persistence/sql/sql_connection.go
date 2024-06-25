// Package sql defines the SQL database management utilities.
package sql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Connection represents a SQL database connection.
type Connection interface {
	// GetDB returns the underlying SQL database connection.
	GetDB() *sqlx.DB
	// GetErrorTranslator returns the error translator associated with the connection.
	GetErrorTranslator() ErrorTranslator
	// GetMigrator returns the database migrator associated with the connection.
	GetMigrator() Migrator
	// GetDialectName returns the name of the SQL dialect used by the connection.
	GetDialectName() string
}

// ConnectionFwOptions defines the options to configure a Connection
type ConnectionFwOptions struct {
	// Postgres defines the PostgresSQL connection information
	Postgres *PostgresInfo `mapstructure:"postgres"`
	// SQLite defines the SQLite connection information
	SQLite *SQLiteInfo
}

// PostgresInfo contains information required to establish a connection to a PostgreSQL database.
type PostgresInfo struct {
	Host     string `mapstructure:"host" valid:"required~hosts is mandatory in Postgres Connection config"`
	Port     int    `mapstructure:"port" valid:"required~port is mandatory in Postgres Connection config"`
	Scheme   string `mapstructure:"scheme" valid:"required~scheme is mandatory in Postgres Connection config"`
	Username string `mapstructure:"username" valid:"required~username is mandatory in Postgres Connection config"`
	Password string `mapstructure:"password" valid:"required~password is mandatory in Postgres Connection config"`
	SSLMode  string `mapstructure:"sslMode" valid:"required~sslMode is mandatory in Postgres Connection config"`
	Database string `mapstructure:"database" valid:"required~database is mandatory in Postgres Connection config"`
}

// SQLiteInfo contains information required to establish a connection to a SQLite database.
type SQLiteInfo struct {
	// ConnectionString to connect to sqlite databases
	ConnectionString string `valid:"required~connectionString is mandatory in SQLite Connection config"`
}

// ConnectionFw represents a database connection framework.
type ConnectionFw struct {
	//db defines the underlying SQL database connection.
	db *sqlx.DB
	// dialect defines the SQL dialect used by the connection.
	dialect Dialect
	// dialectName defines the name of the SQL dialect used by the connection.
	dialectName string
}

// NewConnectionFw returns a new ConnectionFw with the specified ConnectionFwOptions.
func NewConnectionFw(options ConnectionFwOptions) (*ConnectionFw, error) {
	onlyOneNil, dialectName, config := checkOnlyOneNilDialect(options)
	if !onlyOneNil {
		return nil, fmt.Errorf("one and only one db dialect must be not nil")
	}

	d, ok := GetDialect(dialectName)
	if !ok {
		return nil, fmt.Errorf("unsupported dialect [%s]", dialectName)
	}

	connectionData, err := d.GetConnectionData(config)
	if err != nil {
		return nil, err
	}

	db, err := d.Connect(*connectionData)
	if err != nil {
		return nil, err
	}

	return &ConnectionFw{
		dialect:     d,
		db:          db,
		dialectName: dialectName,
	}, nil
}

func (c ConnectionFw) GetDB() *sqlx.DB {
	return c.db
}

func (c ConnectionFw) GetErrorTranslator() ErrorTranslator {
	return c.dialect
}

func (c ConnectionFw) GetMigrator() Migrator {
	return c.dialect
}

func (c ConnectionFw) GetDialectName() string {
	return c.dialectName
}

var _ Connection = (*ConnectionFw)(nil)

func checkOnlyOneNilDialect(options ConnectionFwOptions) (bool, string, interface{}) {
	count := 0
	var propName string
	var config interface{}

	if options.Postgres != nil {
		count++
		propName = "postgres"
		config = options.Postgres
	}
	if options.SQLite != nil {
		count++
		propName = "sqlite"
		config = options.SQLite
	}
	if count == 1 {
		return true, propName, config
	}
	return false, "", nil
}
