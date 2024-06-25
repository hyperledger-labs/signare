package metricrecorder

import (
	"errors"
)

// NewCounterVectorAdapterOptions defines the options to create a new NewCounterVectorAdapter
type NewCounterVectorAdapterOptions struct {
	// Name of the metric
	Name string
	// Help string accompanying the metric
	Help string
	// Array of labels of the metric dimensions
	Labels []string
}

// NewHistogramAdapterOptions defines the options to create a new NewHistogramAdapter
type NewHistogramAdapterOptions struct {
	// Name of the metric
	Name string
	// Help string accompanying the metric
	Help string
	// Array of labels of the metric dimensions
	Labels []string
	// Array of buckets to record measurements
	Buckets []float64
}

// NewGaugeAdapterOptions defines the options to create a new NewGaugeAdapter
type NewGaugeAdapterOptions struct {
	// Name of the metric
	Name string
	// Help string accompanying the metric
	Help string
	// Array of labels of the metric dimensions
	Labels []string
}

// NewSummaryAdapterOptions defines the options to create a new NewSummaryAdapter
type NewSummaryAdapterOptions struct {
	// Name of the metric
	Name string
	// Help string accompanying the metric
	Help string
	// Array of labels of the metric dimensions
	Labels []string
}

var _ MetricRecorder = (*DefaultMetricRecorder)(nil)

type DefaultMetricRecorder struct {
	adapter MetricsRecorderAdapter
}

// DefaultMetricRecorderOptions defines the options to create a new DefaultMetricRecorder
type DefaultMetricRecorderOptions struct {
	MetricsRecorderAdapter MetricsRecorderAdapter
}

// ProvideDefaultMetricRecorder creates a DefaultMetricRecorder with the given DefaultMetricRecorderOptions
func ProvideDefaultMetricRecorder(options DefaultMetricRecorderOptions) (*DefaultMetricRecorder, error) {
	if options.MetricsRecorderAdapter == nil {
		return nil, errors.New("MetricRecorderAdapter is mandatory")
	}
	return &DefaultMetricRecorder{
		adapter: options.MetricsRecorderAdapter,
	}, nil
}

// NewCounterVector creates a CounterVector with the given name, labels and help
func (recorder *DefaultMetricRecorder) NewCounterVector(name string, labels []string, help string) (CounterVector, error) {
	counterAdapter, err := recorder.adapter.NewCounterVector(NewCounterVectorAdapterOptions{
		Name:   name,
		Labels: labels,
		Help:   help,
	})

	if err != nil {
		return nil, err
	}

	return DefaultCounterVector{
		adapter: counterAdapter,
	}, nil
}

// NewHistogramVector creates a HistogramVector with the given name, labels, buckets and help
func (recorder *DefaultMetricRecorder) NewHistogramVector(name string, labels []string, buckets []float64, help string) (HistogramVector, error) {
	histogramAdapter, err := recorder.adapter.NewHistogramVector(NewHistogramAdapterOptions{
		Name:    name,
		Help:    help,
		Labels:  labels,
		Buckets: buckets,
	})

	if err != nil {
		return nil, err
	}

	return DefaultHistogramVector{
		adapter: histogramAdapter,
	}, nil
}

// NewGaugeVector creates a GaugeVector with the given name, labels and help
func (recorder *DefaultMetricRecorder) NewGaugeVector(name string, labels []string, help string) (GaugeVector, error) {
	gaugeAdapter, err := recorder.adapter.NewGauge(NewGaugeAdapterOptions{
		Name:   name,
		Help:   help,
		Labels: labels,
	})

	if err != nil {
		return nil, err
	}

	return DefaultGaugeVector{
		adapter: gaugeAdapter,
	}, nil
}

// NewSummaryVector creates a SummaryVector with the given name, labels and help
func (recorder *DefaultMetricRecorder) NewSummaryVector(name string, labels []string, help string) (SummaryVector, error) {
	summaryAdapter, err := recorder.adapter.NewSummary(NewSummaryAdapterOptions{
		Name:   name,
		Help:   help,
		Labels: labels,
	})

	if err != nil {
		return nil, err
	}

	return DefaultSummaryVector{
		adapter: summaryAdapter,
	}, nil
}
