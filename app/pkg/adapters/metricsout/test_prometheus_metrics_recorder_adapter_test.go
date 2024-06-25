package metricsout_test

import (
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/adapters/metricsout"

	"github.com/stretchr/testify/require"
)

func TestRenderPrometheusMetricDescriptor(t *testing.T) {
	descriptor, err := metricsout.RenderPrometheusMetricDescriptor("test_ns", "test_metric", map[string]string{
		"label1": "value1",
		"label2": "value2",
	})

	require.NoError(t, err)
	expectedDescriptor := `test_ns_test_metric{label1="value1",label2="value2"}`
	require.Equal(t, expectedDescriptor, descriptor)

}
