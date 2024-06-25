package metricsout_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"

	"github.com/stretchr/testify/require"
)

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_HistogramAdapter_Observe(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	histogram, err := metricsRecorderAdapter.NewHistogramVector(metricrecorder.NewHistogramAdapterOptions{
		Name:   "test_histogram",
		Labels: []string{"histogram_label1", "histogram_label2"},
	})
	require.NoError(t, err)

	labels := map[string]string{"histogram_label1": "value1", "histogram_label2": "value2"}

	histogram.Observe(metricrecorder.ObserveHistogramVectorAdapterOptions{Labels: labels, Value: 5})
	histogram.Observe(metricrecorder.ObserveHistogramVectorAdapterOptions{Labels: labels, Value: 5})
	histogram.Observe(metricrecorder.ObserveHistogramVectorAdapterOptions{Labels: labels, Value: 5})

	histogramValues, err := prometheusTestServer.GetHistogramVectorMetricValue(histogram, labels)
	require.NoError(t, err)

	require.Equal(t, "15", histogramValues.Sum)
	require.Equal(t, "3", histogramValues.Count)

}

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_HistogramAdapter_Create_CreateTwiceReturnsSame(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	options := metricrecorder.NewHistogramAdapterOptions{
		Name:   "test_histogram_create_twice",
		Labels: []string{"histogram_label1", "histogram_label2"},
	}
	histogram, err := metricsRecorderAdapter.NewHistogramVector(options)
	require.NoError(t, err)
	require.NotNil(t, histogram)

	histogram, err = metricsRecorderAdapter.NewHistogramVector(options)
	require.NoError(t, err)
	require.NotNil(t, histogram)
}
