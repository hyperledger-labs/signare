package metricsout

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"text/template"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TestMetricsRecorderAdapter implements MetricsRecorderAdapter.
type TestMetricsRecorderAdapter struct {
	PrometheusMetricsRecorderAdapter
}

const DefaultTestMetricsNamespace = "test_metrics_ns"

// NewTestMetricsRecorderAdapter creates a new TestMetricsRecorderAdapter insatnce.
func NewTestMetricsRecorderAdapter() (*TestMetricsRecorderAdapter, error) {
	namespace := PrometheusNamespace(DefaultTestMetricsNamespace)
	prometheusMetricsAdapter, err := ProvidePrometheusMetricsRecorderAdapter(PrometheusMetricsRecorderAdapterOptions{
		Namespace: &namespace,
	})

	if err != nil {
		return nil, err
	}

	return &TestMetricsRecorderAdapter{
		PrometheusMetricsRecorderAdapter: *prometheusMetricsAdapter,
	}, nil
}

// PrometheusTestServer Prometheus server to use for testing.
type PrometheusTestServer struct {
	testRecorder *httptest.ResponseRecorder
}

// NewPrometheusTestServer creates a new PrometheusTestServer.
func NewPrometheusTestServer() *PrometheusTestServer {
	return &PrometheusTestServer{}
}

func (prometheusServer *PrometheusTestServer) Init() {
	prometheusServer.testRecorder = httptest.NewRecorder()
}

func (prometheusServer *PrometheusTestServer) Reset() {
	prometheusServer.testRecorder = httptest.NewRecorder()
}

func (prometheusServer *PrometheusTestServer) Stop() {
	prometheusServer.testRecorder = nil
}

func (prometheusServer PrometheusTestServer) GetCounterVectorMetricValue(prometheusAdapter metricrecorder.CounterVectorAdapter, labels map[string]string) (string, error) {
	instance, ok := prometheusAdapter.(PrometheusCounterVectorAdapter)
	if !ok {
		return "", errors.New("not a prometheus counter vector adapter")
	}

	return prometheusServer.doGetMetricValue(instance.metricDescriptorInfo, labels)
}

func (prometheusServer PrometheusTestServer) GetGaugeMetricValue(prometheusAdapter metricrecorder.GaugeVectorAdapter, labels map[string]string) (string, error) {
	instance, ok := prometheusAdapter.(PrometheusGaugeAdapter)
	if !ok {
		return "", errors.New("not a prometheus gauge adapter")
	}

	return prometheusServer.doGetMetricValue(instance.metricDescriptorInfo, labels)
}

type HistogramVecMetricValues struct {
	Sum   string
	Count string
}

func (prometheusServer PrometheusTestServer) GetHistogramVectorMetricValueFor(namespace string, metricName string, labels map[string]string) (*HistogramVecMetricValues, error) {
	return prometheusServer.doGetHistogramVectorMetricValueInternal(metricDescriptorInfo{
		namespace: namespace,
		name:      metricName,
	}, labels)
}

func (prometheusServer PrometheusTestServer) doGetHistogramVectorMetricValueInternal(metricDescriptorInfo metricDescriptorInfo, labels map[string]string) (*HistogramVecMetricValues, error) {
	name := metricDescriptorInfo.name
	metricDescriptorInfo.name = name + "_sum"
	sum, err := prometheusServer.doGetMetricValue(metricDescriptorInfo, labels)
	if err != nil {
		return nil, err
	}

	metricDescriptorInfo.name = name + "_count"
	count, err := prometheusServer.doGetMetricValue(metricDescriptorInfo, labels)
	if err != nil {
		return nil, err
	}

	return &HistogramVecMetricValues{
		Sum:   sum,
		Count: count,
	}, nil
}

func (prometheusServer PrometheusTestServer) GetHistogramVectorMetricValue(prometheusAdapter metricrecorder.HistogramVectorAdapter, labels map[string]string) (*HistogramVecMetricValues, error) {

	instance, ok := prometheusAdapter.(PrometheusHistogramAdapter)
	if !ok {
		return nil, errors.New("not a prometheus histogram adapter")
	}

	return prometheusServer.doGetHistogramVectorMetricValueInternal(instance.metricDescriptorInfo, labels)
}

type SummaryMetricValues struct {
	Sum   string
	Count string
}

func (prometheusServer PrometheusTestServer) GetSummaryMetricValue(prometheusAdapter metricrecorder.SummaryVectorAdapter, labels map[string]string) (*SummaryMetricValues, error) {

	instance, ok := prometheusAdapter.(PrometheusSummaryAdapter)
	if !ok {
		return nil, errors.New("not a prometheus histogram adapter")
	}

	name := instance.metricDescriptorInfo.name
	instance.metricDescriptorInfo.name = name + "_sum"
	sum, err := prometheusServer.doGetMetricValue(instance.metricDescriptorInfo, labels)
	if err != nil {
		return nil, err
	}

	instance.metricDescriptorInfo.name = name + "_count"
	count, err := prometheusServer.doGetMetricValue(instance.metricDescriptorInfo, labels)
	if err != nil {
		return nil, err
	}

	return &SummaryMetricValues{
		Sum:   sum,
		Count: count,
	}, nil
}

func (prometheusServer PrometheusTestServer) doGetMetricValue(metricDescriptorInfo metricDescriptorInfo, labels map[string]string) (string, error) {
	if prometheusServer.testRecorder == nil {
		return "", errors.New("prometheus test server not initialized")
	}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "/metrics", nil)
	if err != nil {
		return "", err
	}
	promhttp.Handler().ServeHTTP(prometheusServer.testRecorder, req)

	metricDesc, err := RenderPrometheusMetricDescriptor(metricDescriptorInfo.namespace, metricDescriptorInfo.name, labels)
	if err != nil {
		return "", err
	}

	responseBody := prometheusServer.testRecorder.Body.String()

	// LastIndex() is used to deal with gauge kind of metrics where all values are returned.
	// The prometheus server can be improved to return an array of strings for these kind of metrics
	metricIndex := strings.LastIndex(responseBody, metricDesc)
	if metricIndex == -1 {
		return "", errors.New("Metric " + metricDesc + " not found")
	}

	metricEnd := metricIndex + len(metricDesc)
	metricEndingLine := strings.Index(responseBody[metricEnd:], "\n")

	if metricEndingLine != -1 {
		metricEndingLine = metricEnd + metricEndingLine
		return strings.TrimSpace(responseBody[metricEnd:metricEndingLine]), nil
	}
	return strings.TrimSpace(responseBody[metricEnd:]), nil
}

const metricTemplate = `{{.Namespace}}_{{.MetricName}}{{"{"}}{{range $labelName,$labelValue := .Labels}}{{$labelName}}="{{$labelValue}}",{{end}}`

func RenderPrometheusMetricDescriptor(namespace string, metricName string, labels map[string]string) (string, error) {
	tpl, err := template.New("prometheus_metric").Parse(metricTemplate)
	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer
	err = tpl.Execute(&outputBuffer, struct {
		Namespace  string
		MetricName string
		Labels     map[string]string
	}{
		Namespace:  namespace,
		MetricName: metricName,
		Labels:     labels,
	})
	if err != nil {
		return "", err
	}

	prometheusDescriptor := outputBuffer.String()
	lastIndex := strings.LastIndex(prometheusDescriptor, ",")

	prometheusDescriptor = prometheusDescriptor[:lastIndex] + strings.Replace(prometheusDescriptor[lastIndex:], ",", "}", 1)

	return prometheusDescriptor, nil
}
