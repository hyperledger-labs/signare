package metricrecorder

func NewNoMetricsRecorder() *NoMetricsRecorder {
	return &NoMetricsRecorder{}
}

var _ MetricRecorder = (*NoMetricsRecorder)(nil)

// NoMetricsRecorder provides an implementation of MetricRecorder that doesn't collect any metric
// It is provided for those scenarios where metrics are not needed or want to be disabled for whatever reason
type NoMetricsRecorder struct {
}

// NewSummaryVector creates a SummaryVector with the given name and labels
func (n NoMetricsRecorder) NewSummaryVector(_ string, _ []string, _ string) (SummaryVector, error) {
	return NoSummaryVector{}, nil
}

// NewGaugeVector creates a GaugeVector with the given name and labels
func (n NoMetricsRecorder) NewGaugeVector(_ string, _ []string, _ string) (GaugeVector, error) {
	return NoGaugeVector{}, nil
}

// NewCounterVector creates a CounterVector with the given name and labels
func (n NoMetricsRecorder) NewCounterVector(_ string, _ []string, _ string) (CounterVector, error) {
	return NoCounterVector{}, nil
}

// NewHistogramVector creates a HistogramVector with the given name and labels
func (n NoMetricsRecorder) NewHistogramVector(_ string, _ []string, _ []float64, _ string) (HistogramVector, error) {
	return NoHistogramVector{}, nil
}

var _ SummaryVector = (*NoSummaryVector)(nil)

type NoSummaryVector struct {
}

// Observe is a stub that does nothing
func (n NoSummaryVector) Observe(_ map[string]string, _ float64) {
	// Nothing to do
}

// GetSummaryVectorAdapter is a stub that does nothing
func (n NoSummaryVector) GetSummaryVectorAdapter() SummaryVectorAdapter {
	return nil
}

var _ GaugeVector = (*NoGaugeVector)(nil)

type NoGaugeVector struct {
}

// Set is a stub that does nothing
func (n NoGaugeVector) Set(_ map[string]string, _ float64) {
	// Nothing to do
}

// Inc is a stub that does nothing
func (n NoGaugeVector) Inc(_ map[string]string) {
	// Nothing to do
}

// Dec is a stub that does nothing
func (n NoGaugeVector) Dec(_ map[string]string) {
	// Nothing to do
}

// Add is a stub that does nothing
func (n NoGaugeVector) Add(_ map[string]string, _ float64) {
	// Nothing to do
}

// Sub is a stub that does nothing
func (n NoGaugeVector) Sub(_ map[string]string, _ float64) {
	// Nothing to do
}

// GetGaugeVectorAdapter is a stub that does nothing
func (n NoGaugeVector) GetGaugeVectorAdapter() GaugeVectorAdapter {
	return nil
}

var _ CounterVector = (*NoCounterVector)(nil)

type NoCounterVector struct {
}

// Inc is a stub that does nothing
func (n NoCounterVector) Inc(_ map[string]string) {
	// Nothing to do
}

// GetCounterVectorAdapter is a stub that does nothing
func (n NoCounterVector) GetCounterVectorAdapter() CounterVectorAdapter {
	return nil
}

var _ HistogramVector = (*NoHistogramVector)(nil)

type NoHistogramVector struct {
}

// Observe is a stub that does nothing
func (n NoHistogramVector) Observe(_ map[string]string, _ float64) {
	// Nothing to do
}

// GetHistogramVectorAdapter is a stub that does nothing
func (n NoHistogramVector) GetHistogramVectorAdapter() HistogramVectorAdapter {
	return nil
}
