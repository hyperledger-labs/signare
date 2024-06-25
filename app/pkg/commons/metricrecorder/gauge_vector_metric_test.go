package metricrecorder_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func TestDefaultGaugeVector_All_Operations(t *testing.T) {
	prometheusMetricRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	options := metricrecorder.DefaultMetricRecorderOptions{
		MetricsRecorderAdapter: prometheusMetricRecorderAdapter,
	}
	metricRecorder, err := metricrecorder.ProvideDefaultMetricRecorder(options)
	require.NoError(t, err)

	gauge, err := metricRecorder.NewGaugeVector("test_gauge_name", []string{"test_gauge_label_1", "test_gauge_label_2"}, "")
	require.NoError(t, err)

	labels := map[string]string{
		"test_gauge_label_1": "test_gauge_label_1_value",
		"test_gauge_label_2": "test_gauge_label_2_value",
	}

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	gauge.Inc(labels)
	gaugeValue, err := prometheusTestServer.GetGaugeMetricValue(gauge.GetGaugeVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "1", gaugeValue)

	gauge.Dec(labels)
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gauge.GetGaugeVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "0", gaugeValue)

	gauge.Set(labels, 35)
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gauge.GetGaugeVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "35", gaugeValue)

	gauge.Add(labels, 10)
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gauge.GetGaugeVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "45", gaugeValue)

	gauge.Sub(labels, 20)
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gauge.GetGaugeVectorAdapter(), labels)
	require.NoError(t, err)
	require.Equal(t, "25", gaugeValue)
}
