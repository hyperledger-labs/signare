// Package metricshttpinfra provides infrastructure to configure and register the metrics HTTP router
package metricshttpinfra

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hyperledger-labs/signare/app/pkg/infra/httpinfra"
)

// MetricsBasePath base path url for the list of Metrics
const defaultMetricBasePath = "/metrics"

// The number of concurrent HTTP requests is limited to MaxRequestsInFlight. See Prometheus docs.
const defaultMaxRequestsInFlight = 10

// Timeout before response of 503 ServiceUnavailable is served. See Prometheus docs.
const defaultTimeoutInMillis = 30000

// PublishMetricsHTTP publishes metrics
func PublishMetricsHTTP(metricRouter *httpinfra.MetricsHTTPRouter, handler http.Handler, path string) error {
	getMetricsOptions := httpinfra.HandlerMatchOptions{Path: path, Methods: []string{"GET"}}

	return metricRouter.RegisterHandlerFunc(getMetricsOptions, handler)
}

// PrometheusMetricsHTTPOptions configures ProvideMetricsHTTP
type PrometheusMetricsHTTPOptions struct {
	HTTPInfra        *httpinfra.MetricsHTTPRouter
	PrometheusConfig PrometheusConfig
}

// PrometheusMetricsHTTPPublished published HTTP metric routes
type PrometheusMetricsHTTPPublished int

// ProvideMetricsHTTP publishes Metrics
func ProvideMetricsHTTP(options PrometheusMetricsHTTPOptions) (PrometheusMetricsHTTPPublished, error) {
	path := defaultMetricBasePath
	maxRequestsInFlight := defaultMaxRequestsInFlight
	timeoutInMillis := defaultTimeoutInMillis

	if options.PrometheusConfig.Path != nil {
		path = *options.PrometheusConfig.Path
	}
	if options.PrometheusConfig.TimeoutInMillis != nil {
		timeoutInMillis = *options.PrometheusConfig.TimeoutInMillis
	}
	if options.PrometheusConfig.MaxRequestsInFlight != nil {
		maxRequestsInFlight = *options.PrometheusConfig.MaxRequestsInFlight
	}

	handler := promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			ErrorLog:            prometheusLoggerWrapper{},
			MaxRequestsInFlight: maxRequestsInFlight,
			Timeout:             time.Duration(timeoutInMillis) * time.Millisecond,
		})

	err := PublishMetricsHTTP(options.HTTPInfra, handler, path)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
