package metricrecorder_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func TestDefaultMetricRecorder_NewCounterVector(t *testing.T) {
	prometheusMetricRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	options := metricrecorder.DefaultMetricRecorderOptions{
		MetricsRecorderAdapter: prometheusMetricRecorderAdapter,
	}
	metricRecorder, err := metricrecorder.ProvideDefaultMetricRecorder(options)
	require.NoError(t, err)

	counter, err := metricRecorder.NewCounterVector("test_counter_general", []string{"test_counter_general_label_1", "test_counter_general_label_2"}, "")
	require.NoError(t, err)

	labels := map[string]string{
		"test_counter_general_label_1": "test_counter_general_label_1_value",
		"test_counter_general_label_2": "test_counter_general_label_2_value",
	}
	counter.Inc(labels)

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	counterValue, err := prometheusTestServer.GetCounterVectorMetricValue(counter.GetCounterVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "1", counterValue)
}
