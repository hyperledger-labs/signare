package metricrecorder

// HistogramVector defines the functionality of a histogram metric
type HistogramVector interface {
	// Observe increments the Histogram with given labels by 1
	Observe(labels map[string]string, value float64)
	// GetHistogramVectorAdapter returns the associated adapter
	GetHistogramVectorAdapter() HistogramVectorAdapter
}

// HistogramVectorAdapter defines the functionality of a histogram metric adapter
type HistogramVectorAdapter interface {
	// Observe returns the associated adapter
	Observe(options ObserveHistogramVectorAdapterOptions)
}

// ObserveHistogramVectorAdapterOptions defines the options to create a new HistogramVectorAdapter
type ObserveHistogramVectorAdapterOptions struct {
	Labels map[string]string
	Value  float64
}

var _ HistogramVector = (*DefaultHistogramVector)(nil)

type DefaultHistogramVector struct {
	adapter HistogramVectorAdapter
}

// GetHistogramVectorAdapter returns the associated adapter
func (dh DefaultHistogramVector) GetHistogramVectorAdapter() HistogramVectorAdapter {
	return dh.adapter
}

// Observe increments the Histogram with given labels by 1
func (dh DefaultHistogramVector) Observe(labels map[string]string, value float64) {
	dh.adapter.Observe(ObserveHistogramVectorAdapterOptions{
		Labels: labels,
		Value:  value,
	})
}
