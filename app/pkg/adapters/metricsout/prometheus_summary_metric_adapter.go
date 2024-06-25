package metricsout

import (
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus"
)

var _ metricrecorder.SummaryVectorAdapter = new(PrometheusSummaryAdapter)

// PrometheusSummaryAdapter implements SummaryVectorAdapter.
type PrometheusSummaryAdapter struct {
	summary *prometheus.SummaryVec
	metricDescriptorInfo
}

// Observe increments the Summary with given labels by 1
func (pc PrometheusSummaryAdapter) Observe(options metricrecorder.ObserveSummaryAdapterOptions) {
	pc.summary.With(options.Labels).Observe(options.Value)
}
