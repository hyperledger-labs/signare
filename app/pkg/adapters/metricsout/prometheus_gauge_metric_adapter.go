package metricsout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus"
)

var _ metricrecorder.GaugeVectorAdapter = new(PrometheusGaugeAdapter)

// PrometheusGaugeAdapter implements GaugeVectorAdapter.
type PrometheusGaugeAdapter struct {
	gauge *prometheus.GaugeVec
	metricDescriptorInfo
}

// Set sets the Gauge to an arbitrary value with given labels
func (p PrometheusGaugeAdapter) Set(options metricrecorder.GaugeVectorArbitraryValueOptions) {
	p.gauge.With(options.Labels).Set(options.Value)
}

// Inc increments the Gauge with given labels by 1. Use Add to increment it by arbitrary values
func (p PrometheusGaugeAdapter) Inc(options metricrecorder.GaugeVectorOptions) {
	p.gauge.With(options.Labels).Inc()
}

// Dec decrements the Gauge with given labels by 1. Use Sub to decrement it by arbitrary values
func (p PrometheusGaugeAdapter) Dec(options metricrecorder.GaugeVectorOptions) {
	p.gauge.With(options.Labels).Dec()
}

// Add adds the given value to the Gauge with provided labels
// The value can be negative, resulting in a decrease of the Gauge
func (p PrometheusGaugeAdapter) Add(options metricrecorder.GaugeVectorArbitraryValueOptions) {
	p.gauge.With(options.Labels).Add(options.Value)
}

// Sub subtracts the given value from the Gauge with provided labels
// The value can be negative, resulting in an increase of the Gauge
func (p PrometheusGaugeAdapter) Sub(options metricrecorder.GaugeVectorArbitraryValueOptions) {
	p.gauge.With(options.Labels).Sub(options.Value)
}
