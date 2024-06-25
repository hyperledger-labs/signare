package metricrecorder

// MetricRecorder defines the functionality of a metric recorder
type MetricRecorder interface {
	// NewCounterVector creates a CounterVector with the given name, labels and help
	NewCounterVector(name string, labels []string, help string) (CounterVector, error)
	// NewHistogramVector creates a HistogramVector with the given name, labels and help
	NewHistogramVector(name string, labels []string, buckets []float64, help string) (HistogramVector, error)
	// NewGaugeVector creates a GaugeVector with the given name, labels and help
	NewGaugeVector(name string, labels []string, help string) (GaugeVector, error)
	// NewSummaryVector creates a SummaryVector with the given name, labels and help
	NewSummaryVector(name string, labels []string, help string) (SummaryVector, error)
}

// MetricsRecorderAdapter defines the functionality of a metric recorder adapter
type MetricsRecorderAdapter interface {
	// NewCounterVector creates a CounterVectorAdapter with the given NewCounterVectorAdapterOptions
	NewCounterVector(options NewCounterVectorAdapterOptions) (CounterVectorAdapter, error)
	// NewHistogramVector creates a HistogramVectorAdapter with the given NewHistogramAdapterOptions
	NewHistogramVector(options NewHistogramAdapterOptions) (HistogramVectorAdapter, error)
	// NewGauge creates a GaugeVectorAdapter with the given NewGaugeAdapterOptions
	NewGauge(options NewGaugeAdapterOptions) (GaugeVectorAdapter, error)
	// NewSummary creates a SummaryVectorAdapter with the given NewSummaryAdapterOptions
	NewSummary(options NewSummaryAdapterOptions) (SummaryVectorAdapter, error)
}
