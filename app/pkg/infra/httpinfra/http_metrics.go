package httpinfra

import (
	"context"
	"errors"

	"github.com/hyperledger-labs/signare/app/pkg/commons/metricrecorder"
	"github.com/hyperledger-labs/signare/app/pkg/infra/requestcontext"
)

// HTTPMetrics is a set of metrics to monitor HTTP API behavior
type HTTPMetrics interface {
	// IncrementForbiddenAccessCounter increments the count of forbidden access attempts
	IncrementForbiddenAccessCounter(ctx context.Context)
}

// IncrementForbiddenAccessCounter increments the count of forbidden access attempts
func (useCase DefaultHTTPMetrics) IncrementForbiddenAccessCounter(ctx context.Context) {
	actionID := defaultActionID
	action, _ := requestcontext.ActionFromContext(ctx)
	if action != nil && *action != "" {
		actionID = *action
	}

	if useCase.forbiddenAccessCounter != nil {
		useCase.forbiddenAccessCounter.Inc(map[string]string{
			"error":  "403",
			"action": actionID,
		})
	}
}

// At this point there's no real defaultActionID
const defaultActionID = "action.yet.undefined"

// DefaultHTTPMetrics are the default set of metrics for the HTTP API
type DefaultHTTPMetrics struct {
	// forbiddenAccessCounter counts the number of not authorized errors
	forbiddenAccessCounter metricrecorder.CounterVector
}

// DefaultHTTPMetricsOptions are the options needed to create a DefaultHTTPMetrics
type DefaultHTTPMetricsOptions struct {
	// MetricRecorder defines the functionality of a metric recorder
	MetricRecorder metricrecorder.MetricRecorder
}

// ProvideDefaultHTTPMetrics create a new DefaultHTTPMetrics with the provided options
func ProvideDefaultHTTPMetrics(options DefaultHTTPMetricsOptions) (*DefaultHTTPMetrics, error) {
	if options.MetricRecorder == nil {
		return nil, errors.New("mandatory 'DefaultRPCInfraResponseHandler' not provided")
	}

	var forbiddenAccessCounter metricrecorder.CounterVector
	if options.MetricRecorder != nil {
		counter, err := options.MetricRecorder.NewCounterVector("forbidden_access_count", []string{"error", "action"}, "total number of not authorized calls to the API")
		if err != nil {
			return nil, err
		}
		forbiddenAccessCounter = counter
	}

	return &DefaultHTTPMetrics{
		forbiddenAccessCounter: forbiddenAccessCounter,
	}, nil
}
