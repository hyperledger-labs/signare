package metricrecorder

// SummaryVector defines the functionality of a summary metric
type SummaryVector interface {
	// Observe increments the Summary with given labels by 1
	Observe(labels map[string]string, value float64)
	// GetSummaryVectorAdapter returns the associated adapter
	GetSummaryVectorAdapter() SummaryVectorAdapter
}

// SummaryVectorAdapter defines the functionality of a summary metric
type SummaryVectorAdapter interface {
	// Observe increments the Summary with given labels by 1
	Observe(options ObserveSummaryAdapterOptions)
}

// ObserveSummaryAdapterOptions options to create a new ObserveSummaryAdapter
type ObserveSummaryAdapterOptions struct {
	Labels map[string]string
	Value  float64
}

var _ SummaryVector = (*DefaultSummaryVector)(nil)

type DefaultSummaryVector struct {
	adapter SummaryVectorAdapter
}

// GetSummaryVectorAdapter returns the associated adapter
func (ds DefaultSummaryVector) GetSummaryVectorAdapter() SummaryVectorAdapter {
	return ds.adapter
}

// Observe increments the Summary with given labels by 1
func (ds DefaultSummaryVector) Observe(labels map[string]string, value float64) {
	ds.adapter.Observe(ObserveSummaryAdapterOptions{
		Labels: labels,
		Value:  value,
	})
}
