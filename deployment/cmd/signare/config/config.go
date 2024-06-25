// Package config defines configuration parameters of the signare.
package config

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

const (
	staticConfigurationFileName      = "signare-config.yml"
	staticConfigurationFileExtension = "yaml"
)

// StaticConfiguration configures the signare system
type StaticConfiguration struct {
	// Logger configuration for server logging.
	Logger *Logger `mapstructure:"logger" valid:"optional"`
	// DatabaseInfo database information.
	DatabaseInfo DatabaseInfo `mapstructure:"database" valid:"required"`
	// RequestContext defines the context of a request
	RequestContext *RequestContext `mapstructure:"requestContext" valid:"optional"`
	// MetricsConfig provides configuration to expose numeric metrics.
	MetricsConfig *MetricsConfig `mapstructure:"metrics" valid:"optional"`
	// HSMModules provides the configuration of the hardware security modules.
	HSMModules HSMModules `mapstructure:"hsmmodules" valid:"required"`
}

// Logger specification
type Logger struct {
	// LogLevel as level of logs to display
	LogLevel string `mapstructure:"logLevel" valid:"required"`
}

// DatabaseInfo configures signare database access
type DatabaseInfo struct {
	// PostgreSQL database configuration
	PostgreSQL *PostgreSQLInfo `mapstructure:"postgresql"`
}

// RequestContext
type RequestContext struct {
	// UserRequestHeader user in request header
	UserRequestHeader string `mapstructure:"userRequestHeader"`
	// ApplicationRequestHeader application in request header
	ApplicationRequestHeader string `mapstructure:"applicationRequestHeader"`
}

// PostgreSQLInfo defines the access to a SQL-compatible database system
type PostgreSQLInfo struct {
	// Host of database system
	Host string `mapstructure:"host" valid:"required~hosts is mandatory in SQL DB config"`
	// Port of database system
	Port int `mapstructure:"port" valid:"required~port is mandatory in SQL DB config"`
	// Scheme of database system
	Scheme string `mapstructure:"scheme" valid:"required~scheme is mandatory in SQL DB config"`
	// Username to use in database system
	Username string `mapstructure:"username" json:"-"`
	// Password to use with username in database system
	Password string `mapstructure:"password" json:"-"`
	// SSLMode to use in database system
	SSLMode string `mapstructure:"sslmode" valid:"required~sslmode is mandatory in SQL DB config"`
	// Database to access to in database system
	Database string `mapstructure:"database" valid:"required~database is mandatory in SQL DB config"`
	// SQLClient database client configuration
	SQLClient *PostgreSQLClient `mapstructure:"sqlClient"`
}

// PostgreSQLClient configures the database client
type PostgreSQLClient struct {
	// MaxIdleConnections max idle connections for the database/sql handle
	MaxIdleConnections *int `mapstructure:"maxIdleConnections"`
	// MaxOpenConnections max open connections for the database/sql handle
	MaxOpenConnections *int `mapstructure:"maxOpenConnections"`
	// MaxConnectionLifetime max connection lifetime for the database/sql handle
	MaxConnectionLifetime *int `mapstructure:"maxConnectionLifetime"`
}

// MetricsConfig configures signare to export metrics
type MetricsConfig struct {
	// Prometheus Metric Record configuration
	PrometheusMetricsConfig *PrometheusMetricsConfig `mapstructure:"prometheus"  valid:"required"`
}

// PrometheusMetricsConfig provides configuration to expose prometheus metrics
type PrometheusMetricsConfig struct {
	// Port where prometheus metrics will be exposed. Default 9780 aligned with not used port from https://github.com/prometheus/prometheus/wiki/Default-port-allocations
	Port *int `mapstructure:"port" valid:"optional"`
	// Path where prometheus
	Path *string `mapstructure:"path" valid:"optional"`
	// The number of concurrent HTTP requests is limited to MaxRequestsInFlight. See golang prometheus client for deeper documentation
	MaxRequestsInFlight *int `mapstructure:"maxRequestsInFlight" valid:"optional"`
	// If handling a request takes longer than Timeout, it is responded to
	// with 503 ServiceUnavailable and a suitable Message. See golang prometheus client for deeper documentation
	TimeoutInMillis *int `mapstructure:"timeoutInMillis" valid:"optional"`
	// Namespace to prefix metric names
	Namespace *string `mapstructure:"namespace" valid:"optional"`
}

// HSMModules configures the hardware security modules.
type HSMModules struct {
	// SoftHSM configuration for SoftHSM.
	SoftHSM *SoftHSMConfig `mapstructure:"softhsm" valid:"optional"`
}

// SoftHSMConfig configures a SoftHSM in the signare.
type SoftHSMConfig struct {
	Library string `mapstructure:"lib" valid:"required"`
}

func GetStaticConfiguration(path string) (*StaticConfiguration, error) {
	viper.SetConfigName(staticConfigurationFileName)
	viper.SetConfigType(staticConfigurationFileExtension)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading static config file placed in path [%s] [error:%w]", path, err)
	}

	var staticConfiguration StaticConfiguration
	err = viper.Unmarshal(&staticConfiguration)
	if err != nil {
		panic(err)
	}
	valid, err := govalidator.ValidateStruct(staticConfiguration)
	if !valid || err != nil {
		return nil, fmt.Errorf("error validating static config file placed in path[%s] [error:%w]", path, err)
	}

	return &staticConfiguration, nil
}
