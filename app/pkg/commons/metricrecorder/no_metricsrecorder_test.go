package metricrecorder_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func Test_All(t *testing.T) {
	noMetricsRecorder := metricrecorder.NewNoMetricsRecorder()

	// CounterVector
	counterVector, err := noMetricsRecorder.NewCounterVector("", nil, "")
	require.NoError(t, err)

	counterVector.Inc(nil)
	counterVectorAdapter := counterVector.GetCounterVectorAdapter()
	require.Nil(t, counterVectorAdapter)

	// GaugeVector
	gaugeVector, err := noMetricsRecorder.NewGaugeVector("", nil, "")
	require.NoError(t, err)

	gaugeVector.Inc(nil)
	gaugeVector.Dec(nil)
	gaugeVector.Set(nil, -1)
	gaugeVector.Add(nil, -1)
	gaugeVector.Sub(nil, -1)
	gaugeVectorAdapter := gaugeVector.GetGaugeVectorAdapter()
	require.Nil(t, gaugeVectorAdapter)

	// HistogramVector
	histogramVector, err := noMetricsRecorder.NewHistogramVector("", nil, []float64{}, "")
	require.NoError(t, err)

	histogramVector.Observe(nil, -1)
	histogramVectorAdapter := histogramVector.GetHistogramVectorAdapter()
	require.Nil(t, histogramVectorAdapter)

	// SummaryVector
	summaryVector, err := noMetricsRecorder.NewSummaryVector("", nil, "")
	require.NoError(t, err)

	summaryVector.Observe(nil, -1)
	summaryVectorAdapter := summaryVector.GetSummaryVectorAdapter()
	require.Nil(t, summaryVectorAdapter)
}
