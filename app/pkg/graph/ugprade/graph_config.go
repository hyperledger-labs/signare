package upgrade

// Config of the application
type Config struct {
	// Libraries configuration
	Libraries LibrariesConfig `valid:"required"`
}

// LibrariesConfig configuration of the different libraries used by the signare
type LibrariesConfig struct {
	// PersistenceFw persistence framework configuration
	PersistenceFw PersistenceFwConfig `valid:"required"`
}

// PersistenceFwConfig persistence framework configuration
type PersistenceFwConfig struct {
	// PostgreSQL configuration to connect to a PostgreSQL database
	PostgreSQL *PostgresSQLConfig `valid:"optional"`
	// SQLite configuration to connect to a SQLite database. SQLite must be used just for testing purposes and not in a production environment.
	SQLite *SQLiteConfig `valid:"optional"`
}

// PostgresSQLConfig configuration to connect to a PostgreSQL database
type PostgresSQLConfig struct {
	// Host of database system
	Host string `valid:"required"`
	// Port of database system. Default value is 5432
	Port *int `valid:"optional"`
	// Scheme of database system. Default value is "postgres"
	Scheme *string `valid:"optional"`
	// Username to use in database system
	Username string `valid:"required"`
	// Password to use with username in database system
	Password string `valid:"required"`
	// SSLMode to use in database system. Default value is "disable", however, it is advised to enable SSL for security reasons
	SSLMode string `valid:"optional"`
	// Database to access to in the database system
	Database string `valid:"required"`
}

// SQLiteConfig configuration for the SQLite client
type SQLiteConfig struct{}
