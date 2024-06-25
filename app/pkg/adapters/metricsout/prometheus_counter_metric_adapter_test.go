package metricsout_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"
	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/stretchr/testify/require"
)

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_CounterAdapter_Increase(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	counter, err := metricsRecorderAdapter.NewCounterVector(metricrecorder.NewCounterVectorAdapterOptions{
		Name:   "test_counter",
		Labels: []string{"label1", "label2"},
	})
	require.NoError(t, err)

	labels := map[string]string{"label1": "value1", "label2": "value2"}

	counter.Inc(metricrecorder.IncreaseCounterAdapterOptions{LabelsValues: labels})
	counter.Inc(metricrecorder.IncreaseCounterAdapterOptions{LabelsValues: labels})

	prometheusTestServer := metricsout.NewPrometheusTestServer()
	prometheusTestServer.Init()
	defer prometheusTestServer.Reset()

	counterValue, err := prometheusTestServer.GetCounterVectorMetricValue(counter, labels)
	require.NoError(t, err)

	require.Equal(t, "2", counterValue)

	otherCounterLabels := map[string]string{"label1": "value3", "label2": "value2"}
	counter.Inc(metricrecorder.IncreaseCounterAdapterOptions{LabelsValues: otherCounterLabels})
	counterValue, err = prometheusTestServer.GetCounterVectorMetricValue(counter, otherCounterLabels)
	require.NoError(t, err)

	require.Equal(t, "1", counterValue)
}

func TestMetricsRecorderOut_PrometheusMetricsRecorderAdapter_CounterAdapter_CreateTwiceReturnsSame(t *testing.T) {
	metricsRecorderAdapter, err := metricsout.NewTestMetricsRecorderAdapter()
	require.NoError(t, err)

	counterOptions := metricrecorder.NewCounterVectorAdapterOptions{
		Name:   "test_counter_create_twice",
		Labels: []string{"label1", "label2"},
	}

	counter, err := metricsRecorderAdapter.NewCounterVector(counterOptions)
	require.NoError(t, err)
	require.NotNil(t, counter)

	counter, err = metricsRecorderAdapter.NewCounterVector(counterOptions)
	require.NoError(t, err)
	require.NotNil(t, counter)
}
