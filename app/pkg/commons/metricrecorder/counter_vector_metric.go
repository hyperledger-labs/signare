// Package metricrecorder defines the utilities to define and manage Prometheus metrics.
package metricrecorder

// CounterVector defines the functionality of a counter metric
type CounterVector interface {
	// Inc increases counter value
	Inc(labelsValues map[string]string)
	// GetCounterVectorAdapter returns the associated adapter
	GetCounterVectorAdapter() CounterVectorAdapter
}

// CounterVectorAdapter defines the functionality of a counter metric adapter
type CounterVectorAdapter interface {
	// Inc increases counter value
	Inc(options IncreaseCounterAdapterOptions)
}

// IncreaseCounterAdapterOptions defines the options to create a new IncreaseCounterAdapter
type IncreaseCounterAdapterOptions struct {
	LabelsValues map[string]string
}

var _ CounterVector = (*DefaultCounterVector)(nil)

type DefaultCounterVector struct {
	adapter CounterVectorAdapter
}

// GetCounterVectorAdapter returns the associated adapter
func (defaultCounter DefaultCounterVector) GetCounterVectorAdapter() CounterVectorAdapter {
	return defaultCounter.adapter
}

// Inc increases counter value
func (defaultCounter DefaultCounterVector) Inc(labelsValues map[string]string) {
	defaultCounter.adapter.Inc(IncreaseCounterAdapterOptions{LabelsValues: labelsValues})
}
