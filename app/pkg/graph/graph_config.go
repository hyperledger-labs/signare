package graph

// Config of the application
type Config struct {
	// BuildConfig build configuration
	BuildConfig *BuildConfig `valid:"optional"`
	// Libraries configuration
	Libraries LibrariesConfig `valid:"required"`
	// RequestContextConfig configure the headers in a request
	RequestContextConfig *RequestContextConfig `valid:"optional"`
}

// BuildConfig defines the information of the current signare build
type BuildConfig struct {
	// BuildTime timestamp on current build
	BuildTime *string `valid:"optional"`
	// Tag in current commit
	Tag *string `valid:"optional"`
	// CommitHash Git commit hash
	CommitHash *string `valid:"optional"`
}

// LibrariesConfig configuration of the different libraries used by the signare
type LibrariesConfig struct {
	// Logger configuration
	Logger *LoggerConfig `valid:"optional"`
	// Metrics configuration
	Metrics *MetricsConfig `valid:"optional"`
	// PersistenceFw persistence framework configuration
	PersistenceFw PersistenceFwConfig `valid:"required"`
	// HSMModules provides the configuration of the hardware security modules.
	HSMModules HSMModules `mapstructure:"hsmmodules" valid:"required"`
}

// LoggerConfig configuration of the logger
type LoggerConfig struct {
	// LogLevel to use for logging. Default level is INFO
	LogLevel *string `valid:"optional"`
}

// MetricsConfig configuration of the metric recorder
type MetricsConfig struct {
	// Prometheus configuration
	Prometheus PrometheusConfig `valid:"required"`
}

// PrometheusConfig Prometheus configuration for the metric recorder
type PrometheusConfig struct {
	// Port where prometheus metrics will be exposed. Default 9780 aligned with not used port from https://github.com/prometheus/prometheus/wiki/Default-port-allocations
	Port *int `valid:"optional"`
	// Path where prometheus is
	Path *string `valid:"optional"`
	// The number of concurrent HTTP requests is limited to MaxRequestsInFlight. See Golang prometheus client for deeper documentation
	MaxRequestsInFlight *int `valid:"optional"`
	// If handling a request takes longer than timeout, the response is 503 ServiceUnavailable and a message. See Golang Prometheus client for deeper documentation
	TimeoutInMillis *int `valid:"optional"`
	// Namespace to prefix metric names. If empty, the default namespace will be applied as defined in package metricsout (defaultNamespace)
	Namespace *string `valid:"optional"`
}

// PersistenceFwConfig persistence framework configuration
type PersistenceFwConfig struct {
	// PostgreSQL configuration to connect to a PostgreSQL database
	PostgreSQL *PostgresSQLConfig `valid:"optional"`
	// SQLite configuration to connect to a SQLite database. SQLite must be used just for testing purposes and not in a production environment.
	SQLite *SQLiteConfig `valid:"optional"`
}

// HSMModules configures the hardware security modules.
type HSMModules struct {
	// SoftHSM configuration for SoftHSM.
	SoftHSM *SoftHSMConfig `mapstructure:"softhsm" valid:"optional"`
}

// SoftHSMConfig configures a SoftHSM.
type SoftHSMConfig struct {
	Library string `mapstructure:"lib" valid:"required"`
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
	// SQLClient database client configuration
	SQLClient *PostgresSQLClientConfig `valid:"optional"`
}

// PostgresSQLClientConfig configuration for the PostgreSQL client
type PostgresSQLClientConfig struct {
	// MaxIdleConnections max idle connections for the database/sql handle
	MaxIdleConnections *int `valid:"optional"`
	// MaxOpenConnections max open connections for the database/sql handle
	MaxOpenConnections *int `valid:"optional"`
	// MaxConnectionLifetime max connection lifetime for the database/sql handle
	MaxConnectionLifetime *int `valid:"optional"`
}

// SQLiteConfig configuration for the SQLite client
type SQLiteConfig struct{}

// RequestContextConfig configures the keys in the headers of a request
type RequestContextConfig struct {
	// UserHeaderKey is the header key to define the user of a request
	UserHeaderKey string `mapstructure:"userHeaderKey"`
	// ApplicationHeaderKey is the header key to define the application of a request
	ApplicationHeaderKey string `mapstructure:"applicationHeaderKey"`
}
