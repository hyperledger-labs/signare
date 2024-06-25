package metricsout

import (
	"errors"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus"
)

var _ metricrecorder.MetricsRecorderAdapter = new(PrometheusMetricsRecorderAdapter)

// Prometheus namespace should have a (single-word) application prefix relevant to the domain the metric belongs to.
// https://prometheus.io/docs/practices/naming/
const defaultNamespace PrometheusNamespace = "signer"

// PrometheusNamespace namespace to use in Prometheus for the Signare.
type PrometheusNamespace string

// NewCounterVector creates a CounterVectorAdapter with the given NewCounterVectorAdapterOptions
func (p PrometheusMetricsRecorderAdapter) NewCounterVector(options metricrecorder.NewCounterVectorAdapterOptions) (metricrecorder.CounterVectorAdapter, error) {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: string(p.namespace),
			Name:      options.Name,
			Help:      options.Help,
		},
		options.Labels,
	)

	err := prometheus.Register(counter)
	if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
		return nil, err
	}

	return PrometheusCounterVectorAdapter{
		counter: counter,
		metricDescriptorInfo: metricDescriptorInfo{
			namespace: string(p.namespace),
			name:      options.Name,
		},
	}, nil
}

// NewHistogramVector creates a HistogramVectorAdapter with the given NewHistogramAdapterOptions
func (p PrometheusMetricsRecorderAdapter) NewHistogramVector(options metricrecorder.NewHistogramAdapterOptions) (metricrecorder.HistogramVectorAdapter, error) {
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: string(p.namespace),
			Name:      options.Name,
			Help:      options.Help,
			Buckets:   options.Buckets,
		}, options.Labels)

	err := prometheus.Register(histogram)
	if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
		return nil, err
	}

	return PrometheusHistogramAdapter{
		histogram: histogram,
		metricDescriptorInfo: metricDescriptorInfo{
			namespace: string(p.namespace),
			name:      options.Name,
		},
	}, nil
}

// NewGauge creates a GaugeVectorAdapter with the given NewGaugeAdapterOptions
func (p PrometheusMetricsRecorderAdapter) NewGauge(options metricrecorder.NewGaugeAdapterOptions) (metricrecorder.GaugeVectorAdapter, error) {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: string(p.namespace),
			Name:      options.Name,
			Help:      options.Help,
		}, options.Labels)

	err := prometheus.Register(gauge)
	if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
		return nil, err
	}

	return PrometheusGaugeAdapter{
		gauge: gauge,
		metricDescriptorInfo: metricDescriptorInfo{
			namespace: string(p.namespace),
			name:      options.Name,
		},
	}, nil
}

// NewSummary creates a SummaryVectorAdapter with the given NewSummaryAdapterOptions
func (p PrometheusMetricsRecorderAdapter) NewSummary(options metricrecorder.NewSummaryAdapterOptions) (metricrecorder.SummaryVectorAdapter, error) {
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: string(p.namespace),
			Name:      options.Name,
			Help:      options.Help,
		}, options.Labels)

	err := prometheus.Register(summary)
	if err != nil && !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
		return nil, err
	}

	return PrometheusSummaryAdapter{
		summary: summary,
		metricDescriptorInfo: metricDescriptorInfo{
			namespace: string(p.namespace),
			name:      options.Name,
		},
	}, nil
}

// PrometheusMetricsRecorderAdapter implements MetricsRecorderAdapter.
type PrometheusMetricsRecorderAdapter struct {
	namespace PrometheusNamespace
}

// PrometheusMetricsRecorderAdapterOptions options to create a new PrometheusMetricsRecorderAdapter.
type PrometheusMetricsRecorderAdapterOptions struct {
	Namespace *PrometheusNamespace
}

// ProvidePrometheusMetricsRecorderAdapter creates a PrometheusMetricsRecorderAdapter with the given PrometheusMetricsRecorderAdapterOptions
func ProvidePrometheusMetricsRecorderAdapter(options PrometheusMetricsRecorderAdapterOptions) (*PrometheusMetricsRecorderAdapter, error) {
	namespace := defaultNamespace
	if options.Namespace != nil {
		namespace = *options.Namespace
	}
	return &PrometheusMetricsRecorderAdapter{
		namespace: namespace,
	}, nil
}
