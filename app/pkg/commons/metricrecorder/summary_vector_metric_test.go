package metricrecorder_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func TestDefaultSummaryVector_Observe(t *testing.T) {
	prometheusMetricRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	options := metricrecorder.DefaultMetricRecorderOptions{
		MetricsRecorderAdapter: prometheusMetricRecorderAdapter,
	}
	metricRecorder, err := metricrecorder.ProvideDefaultMetricRecorder(options)
	require.NoError(t, err)

	summary, err := metricRecorder.NewSummaryVector("test_summary_name", []string{"test_summary_label_1", "test_summary_label_2"}, "")
	require.NoError(t, err)

	labels := map[string]string{
		"test_summary_label_1": "test_summary_label_1_value",
		"test_summary_label_2": "test_summary_label_2_value",
	}

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	summary.Observe(labels, 5)
	summary.Observe(labels, 100)
	summaryValue, err := prometheusTestServer.GetSummaryMetricValue(summary.GetSummaryVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "2", summaryValue.Count)
	require.Equal(t, "105", summaryValue.Sum)
}
