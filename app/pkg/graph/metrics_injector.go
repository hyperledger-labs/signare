//go:build wireinject

package graph

import (
	"github.com/google/wire"
	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
)

type prometheusMetricsGraph struct {
	metricRecorderAdapter metricrecorder.MetricsRecorderAdapter
	metricRecorder        metricrecorder.MetricRecorder
}

var prometheusMetricsSet = wire.NewSet(
	wire.Struct(new(prometheusMetricsGraph), "*"),

	// Metric Recorder
	metricrecorder.ProvideDefaultMetricRecorder,
	wire.Bind(new(metricrecorder.MetricRecorder), new(*metricrecorder.DefaultMetricRecorder)),
	wire.Struct(new(metricrecorder.DefaultMetricRecorderOptions), "MetricsRecorderAdapter"),

	// Prometheus Metrics Recorder Handler
	metricsout.ProvidePrometheusMetricsRecorderAdapter,
	providePrometheusNamespace,
	wire.Struct(new(metricsout.PrometheusMetricsRecorderAdapterOptions), "*"),
	wire.Bind(new(metricrecorder.MetricsRecorderAdapter), new(*metricsout.PrometheusMetricsRecorderAdapter)),
)

func initializePrometheusMetrics(config Config) (*prometheusMetricsGraph, error) {
	wire.Build(prometheusMetricsSet)
	return &prometheusMetricsGraph{}, nil
}

func providePrometheusNamespace(config Config) *metricsout.PrometheusNamespace {
	if config.Libraries.Metrics != nil && config.Libraries.Metrics.Prometheus.Namespace != nil {
		ns := metricsout.PrometheusNamespace(*config.Libraries.Metrics.Prometheus.Namespace)
		return &ns
	}
	return nil
}

type dummyMetricsGraph struct {
	metricRecorder metricrecorder.MetricRecorder
}

var dummyMetricsSet = wire.NewSet(
	wire.Struct(new(dummyMetricsGraph), "*"),

	// Metric Recorder
	metricrecorder.NewNoMetricsRecorder,
	wire.Bind(new(metricrecorder.MetricRecorder), new(*metricrecorder.NoMetricsRecorder)),
)

func InitializeDummyMetrics() (*dummyMetricsGraph, error) {
	wire.Build(dummyMetricsSet)
	return &dummyMetricsGraph{}, nil
}
