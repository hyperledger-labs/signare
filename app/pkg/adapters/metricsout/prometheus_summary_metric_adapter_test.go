package metricsout_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"

	"github.com/stretchr/testify/require"
)

func TestPrometheusSummaryAdapter_Observe(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	summary, err := metricsRecorderAdapter.NewSummary(metricrecorder.NewSummaryAdapterOptions{
		Name:   "test_summary",
		Labels: []string{"summary_label1", "summary_label2"},
	})
	require.NoError(t, err)

	labels := map[string]string{"summary_label1": "value1", "summary_label2": "value2"}

	summary.Observe(metricrecorder.ObserveSummaryAdapterOptions{
		Labels: labels,
		Value:  42,
	})

	summary.Observe(metricrecorder.ObserveSummaryAdapterOptions{
		Labels: labels,
		Value:  84,
	})

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	counterValue, err := prometheusTestServer.GetSummaryMetricValue(summary, labels)
	require.NoError(t, err)

	require.Equal(t, "126", counterValue.Sum)
	require.Equal(t, "2", counterValue.Count)
}

func TestPrometheusSummaryAdapter_Create_CreateTwiceReturnsSame(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	options := metricrecorder.NewSummaryAdapterOptions{
		Name:   "test_summary_create_twice",
		Labels: []string{"summary_label1", "summary_label2"},
	}
	summary, err := metricsRecorderAdapter.NewSummary(options)
	require.NoError(t, err)
	require.NotNil(t, summary)

	summary, err = metricsRecorderAdapter.NewSummary(options)
	require.NoError(t, err)
	require.NotNil(t, summary)
}
