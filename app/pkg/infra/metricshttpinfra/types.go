package metricshttpinfra

import (
	"context"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
)

// PrometheusConfig Prometheus configuration for the metric recorder
type PrometheusConfig struct {
	// Path where prometheus is
	Path *string
	// The number of concurrent HTTP requests is limited to MaxRequestsInFlight. See Golang prometheus client for deeper documentation
	MaxRequestsInFlight *int
	// If handling a request takes longer than timeout, the response is 503 ServiceUnavailable and a message. See Golang Prometheus client for deeper documentation
	TimeoutInMillis *int
}

// prometheusLoggerWrapper defines a wrapper struct for the logger
type prometheusLoggerWrapper struct {
}

// Println prints prometheusLoggerWrapper
func (p prometheusLoggerWrapper) Println(v ...interface{}) {
	logger.LogEntry(context.Background()).Infof("%v ", v...)
}
