package metricsout_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"

	"github.com/stretchr/testify/require"
)

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_GaugeAdapter_All_Operations(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	gaugeMetric, err := metricsRecorderAdapter.NewGauge(metricrecorder.NewGaugeAdapterOptions{
		Name:   "test_gauge",
		Labels: []string{"gauge_label1", "gauge_label2"},
	})
	require.NoError(t, err)

	labels := map[string]string{"gauge_label1": "value1", "gauge_label2": "value2"}

	gaugeMetric.Inc(metricrecorder.GaugeVectorOptions{Labels: labels})
	gaugeMetric.Inc(metricrecorder.GaugeVectorOptions{Labels: labels})

	gaugeValue, err := prometheusTestServer.GetGaugeMetricValue(gaugeMetric, labels)
	require.NoError(t, err)
	require.Equal(t, "2", gaugeValue)

	gaugeMetric.Dec(metricrecorder.GaugeVectorOptions{Labels: labels})
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gaugeMetric, labels)
	require.NoError(t, err)
	require.Equal(t, "1", gaugeValue)

	gaugeMetric.Set(metricrecorder.GaugeVectorArbitraryValueOptions{Labels: labels, Value: 5})
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gaugeMetric, labels)
	require.NoError(t, err)
	require.Equal(t, "5", gaugeValue)

	gaugeMetric.Add(metricrecorder.GaugeVectorArbitraryValueOptions{Labels: labels, Value: 20})
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gaugeMetric, labels)
	require.NoError(t, err)
	require.Equal(t, "25", gaugeValue)

	gaugeMetric.Sub(metricrecorder.GaugeVectorArbitraryValueOptions{Labels: labels, Value: 10})
	gaugeValue, err = prometheusTestServer.GetGaugeMetricValue(gaugeMetric, labels)
	require.NoError(t, err)
	require.Equal(t, "15", gaugeValue)

}

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_GaugeAdapter_CreateTwiceReturnsSame(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	options := metricrecorder.NewGaugeAdapterOptions{
		Name:   "test_gauge_create_twice",
		Labels: []string{"gauge_label1", "gauge_label2"},
	}

	gauge, err := metricsRecorderAdapter.NewGauge(options)
	require.NoError(t, err)
	require.NotNil(t, gauge)

	gauge, err = metricsRecorderAdapter.NewGauge(options)
	require.NoError(t, err)
	require.NotNil(t, gauge)
}
