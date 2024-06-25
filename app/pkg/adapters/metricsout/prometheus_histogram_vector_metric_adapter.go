package metricsout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus"
)

var _ metricrecorder.HistogramVectorAdapter = new(PrometheusHistogramAdapter)

// PrometheusHistogramAdapter implements HistogramVectorAdapter.
type PrometheusHistogramAdapter struct {
	histogram *prometheus.HistogramVec
	metricDescriptorInfo
}

// Observe increments the Histogram with given labels by 1
func (p PrometheusHistogramAdapter) Observe(options metricrecorder.ObserveHistogramVectorAdapterOptions) {
	p.histogram.With(options.Labels).Observe(options.Value)
}
