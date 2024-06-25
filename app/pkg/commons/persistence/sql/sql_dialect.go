package sql

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
)

var dialectsMu sync.RWMutex
var dialects = make(map[string]Dialect)

type MigrationVersion struct {
	Version     int
	Dirty       bool
	Description string
}

type SetVersionInput struct {
	MigrationVersion
}

type GetVersionOutput struct {
	MigrationVersion
}

type Dialect interface {
	Opener
	Migrator
	ErrorTranslator
}

type ConnectionData struct {
	ConnectionString string
	Driver           string
}

// Opener defines functionality to establish a database connection.
type Opener interface {
	// GetConnectionData retrieves connection data from the provided configuration.
	GetConnectionData(config interface{}) (*ConnectionData, error)
	// Connect establishes a database connection using the provided connection data.
	Connect(connectionData ConnectionData) (*sqlx.DB, error)
}

// Migrator manages database migrations.
type Migrator interface {
	// OpenConnection opens a connection to the database.
	OpenConnection(ctx context.Context) error
	// InitMigration initializes the migration process.
	InitMigration(ctx context.Context, migrationsTablePrefix *string) error
	// GetMigrationVersion retrieves the current migration version.
	GetMigrationVersion(ctx context.Context) (*GetVersionOutput, error)
	// SetMigrationVersion sets the migration version.
	SetMigrationVersion(ctx context.Context, input SetVersionInput) error
	// RunMigration runs a migration script.
	RunMigration(ctx context.Context, migration string) error
	// CloseConnection closes the connection to the database.
	CloseConnection(ctx context.Context) error
}

// ErrorTranslator translates SQL errors for better context.
type ErrorTranslator interface {
	// TranslateError translates an SQL error.
	TranslateError(ctx context.Context, originalErr error, defaultError error) error
}

func GetDialect(s string) (Dialect, bool) {
	dialectsMu.RLock()
	d, ok := dialects[s]
	dialectsMu.RUnlock()
	return d, ok
}

func RegisterDialect(name string, dialect Dialect) {
	dialectsMu.Lock()
	defer dialectsMu.Unlock()
	if dialect == nil {
		panic("Register dialect is nil")
	}
	if _, dup := dialects[name]; dup {
		panic("Register called twice for dialect " + name)
	}
	dialects[name] = dialect
}
