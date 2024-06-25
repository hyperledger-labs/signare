// Package metricsout defines the implementation of the output adapter of the Prometheus metrics.
package metricsout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus"
)

var _ metricrecorder.CounterVectorAdapter = new(PrometheusCounterVectorAdapter)

// PrometheusCounterVectorAdapter implements CounterVectorAdapter.
type PrometheusCounterVectorAdapter struct {
	counter *prometheus.CounterVec
	metricDescriptorInfo
}

// Inc increases counter value
func (pc PrometheusCounterVectorAdapter) Inc(options metricrecorder.IncreaseCounterAdapterOptions) {
	pc.counter.With(options.LabelsValues).Inc()
}
