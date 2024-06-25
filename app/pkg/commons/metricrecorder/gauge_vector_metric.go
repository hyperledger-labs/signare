package metricrecorder

// GaugeVector defines the functionality of a gauge metric
type GaugeVector interface {
	// Set sets the Gauge to an arbitrary value with given labels
	Set(labels map[string]string, value float64)
	// Inc increments the Gauge with given labels by 1. Use Add to increment it by arbitrary values
	Inc(labels map[string]string)
	// Dec decrements the Gauge with given labels by 1. Use Sub to decrement it by arbitrary values
	Dec(labels map[string]string)
	// Add adds the given value to the Gauge with provided labels
	// The value can be negative, resulting in a decrease of the Gauge
	Add(labels map[string]string, value float64)
	// Sub subtracts the given value from the Gauge with provided labels
	// The value can be negative, resulting in an increase of the Gauge
	Sub(labels map[string]string, value float64)
	// GetGaugeVectorAdapter returns the associated adapter
	GetGaugeVectorAdapter() GaugeVectorAdapter
}

// GaugeVectorAdapter defines the functionality of a gauge metric adapter
type GaugeVectorAdapter interface {
	// Set sets the Gauge to an arbitrary value with given labels
	Set(options GaugeVectorArbitraryValueOptions)
	// Inc increments the Gauge with given labels by 1. Use Add to increment it by arbitrary values
	Inc(options GaugeVectorOptions)
	// Dec decrements the Gauge with given labels by 1. Use Sub to decrement it by arbitrary values
	Dec(GaugeVectorOptions)
	// Add adds the given value to the Gauge with provided labels
	// The value can be negative, resulting in a decrease of the Gauge
	Add(options GaugeVectorArbitraryValueOptions)
	// Sub subtracts the given value from the Gauge with provided labels
	// The value can be negative, resulting in an increase of the Gauge
	Sub(options GaugeVectorArbitraryValueOptions)
}

// GaugeVectorOptions defines the options to create a new GaugeVector
type GaugeVectorOptions struct {
	Labels map[string]string
}

// GaugeVectorArbitraryValueOptions defines the options to create a new GaugeVectorArbitraryValue
type GaugeVectorArbitraryValueOptions struct {
	Labels map[string]string
	Value  float64
}

var _ GaugeVector = (*DefaultGaugeVector)(nil)

type DefaultGaugeVector struct {
	adapter GaugeVectorAdapter
}

// GetGaugeVectorAdapter returns the associated adapter
func (d DefaultGaugeVector) GetGaugeVectorAdapter() GaugeVectorAdapter {
	return d.adapter
}

// Set sets the Gauge to an arbitrary value with given labels
func (d DefaultGaugeVector) Set(labels map[string]string, value float64) {
	d.adapter.Set(GaugeVectorArbitraryValueOptions{
		Labels: labels,
		Value:  value,
	})
}

// Inc increments the Gauge with given labels by 1. Use Add to increment it by arbitrary values
func (d DefaultGaugeVector) Inc(labels map[string]string) {
	d.adapter.Inc(GaugeVectorOptions{Labels: labels})
}

// Dec decrements the Gauge with given labels by 1. Use Sub to decrement it by arbitrary values
func (d DefaultGaugeVector) Dec(labels map[string]string) {
	d.adapter.Dec(GaugeVectorOptions{Labels: labels})
}

// Add adds the given value to the Gauge with provided labels
// The value can be negative, resulting in a decrease of the Gauge
func (d DefaultGaugeVector) Add(labels map[string]string, value float64) {
	d.adapter.Add(GaugeVectorArbitraryValueOptions{
		Labels: labels,
		Value:  value,
	})
}

// Sub subtracts the given value from the Gauge with provided labels
// The value can be negative, resulting in an increase of the Gauge
func (d DefaultGaugeVector) Sub(labels map[string]string, value float64) {
	d.adapter.Sub(GaugeVectorArbitraryValueOptions{
		Labels: labels,
		Value:  value,
	})
}
