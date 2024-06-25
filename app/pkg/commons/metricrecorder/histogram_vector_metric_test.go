package metricrecorder_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func TestDefaultHistogramVector_Observe(t *testing.T) {
	prometheusMetricRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	options := metricrecorder.DefaultMetricRecorderOptions{
		MetricsRecorderAdapter: prometheusMetricRecorderAdapter,
	}
	metricRecorder, err := metricrecorder.ProvideDefaultMetricRecorder(options)
	require.NoError(t, err)

	histogram, err := metricRecorder.NewHistogramVector("test_histogram_name", []string{"test_histogram_label_1", "test_histogram_label_2"}, []float64{}, "")
	require.NoError(t, err)

	labels := map[string]string{
		"test_histogram_label_1": "test_histogram_label_1_value",
		"test_histogram_label_2": "test_histogram_label_2_value",
	}

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	histogram.Observe(labels, 5)
	histogram.Observe(labels, 20)
	histogramValue, err := prometheusTestServer.GetHistogramVectorMetricValue(histogram.GetHistogramVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "2", histogramValue.Count)
	require.Equal(t, "25", histogramValue.Sum)
}
