// nolint: loggercheck
package logger

import (
	"context"
	"log/slog"
	"time"
)

// Tracer represents a logger with tracing capabilities. It allows the creation of a scoped tracing object with a specific namespace.
type Tracer interface {
	// AddProperty allows the inclusion of permanent properties to the Tracer's namespace object.
	AddProperty(key string, value any)
	// Trace creates a log entry with LevelInfo containing the tracer's object.
	Trace(msg string)
	// TraceWithData creates a log entry with LevelInfo containing the tracer's object.
	// The data arguments are placed inside the tracer's namespace object.
	TraceWithData(msg string, data map[string]any)
	// Warn creates a log entry with LevelWarn containing the tracer's object.
	Warn(msg string)
	// WarnWithData creates a log entry with LevelWarn containing the tracer's object.
	// The data arguments are placed inside the tracer's namespace object.
	WarnWithData(msg string, data map[string]any)
	// Error creates a log entry with LevelError containing the tracer's object.
	// The message is placed outside the tracer's namespace object.
	Error(msg string)
	// Errorf creates a log entry with LevelError containing the tracer's object,
	// formatting the message using a format string and arguments provided.
	Errorf(format string, data ...interface{})
	// ErrorWithData creates a log entry with LevelError containing the tracer's object.
	// The message is placed outside the tracer's namespace object, while the data arguments are placed inside of it.
	ErrorWithData(msg string, data map[string]any)
	// Debug creates a log entry with LevelDebug containing the tracer's object.
	Debug(msg string)
	// Debugf creates a log entry with LevelDebug containing the tracer's object,
	// formatting the message using a format string and arguments provided.
	Debugf(format string, data ...interface{})
	// DebugWithData creates a log entry with LevelDebug containing the tracer's object.
	// The data arguments are placed inside the tracer's namespace object.
	DebugWithData(msg string, data map[string]any)
}

type tracer struct {
	logger     Logger
	namespace  string
	attributes []slog.Attr
}

var _ Tracer = new(tracer)

// NewTracer returns a new Tracer for the specified namespace.
func NewTracer(ctx context.Context) Tracer {
	return &tracer{
		logger:     LogEntry(ctx),
		namespace:  "properties", // the dynamic part of the structured log will be contained in a single attribute named "properties"
		attributes: make([]slog.Attr, 0),
	}
}

func (t *tracer) getGroup() slog.Attr {
	return slog.Attr{
		Key:   t.namespace,
		Value: slog.GroupValue(t.attributes...),
	}
}

func (t *tracer) AddProperty(key string, value any) {
	switch v := value.(type) {
	case string:
		t.attributes = append(t.attributes, slog.String(key, v))
		return
	case int64:
		t.attributes = append(t.attributes, slog.Int64(key, v))
		return
	case int:
		t.attributes = append(t.attributes, slog.Int(key, v))
		return
	case uint64:
		t.attributes = append(t.attributes, slog.Uint64(key, v))
		return
	case float64:
		t.attributes = append(t.attributes, slog.Float64(key, v))
		return
	case bool:
		t.attributes = append(t.attributes, slog.Bool(key, v))
		return
	case time.Time:
		t.attributes = append(t.attributes, slog.Time(key, v))
		return
	case time.Duration:
		t.attributes = append(t.attributes, slog.Duration(key, v))
		return
	}
	t.attributes = append(t.attributes, slog.Attr{
		Key:   key,
		Value: slog.AnyValue(value),
	})
}

func (t *tracer) Trace(msg string) {
	t.logger.WithArguments(t.getGroup()).Info(msg)
}

func (t *tracer) TraceWithData(msg string, data map[string]any) {
	tracerCopy := *t
	for k, v := range data {
		tracerCopy.AddProperty(k, v)
	}
	t.logger.WithArguments(tracerCopy.getGroup()).Info(msg)
}

func (t *tracer) Warn(msg string) {
	t.logger.WithArguments(t.getGroup()).Warn(msg)
}

func (t *tracer) WarnWithData(msg string, data map[string]any) {
	tracerCopy := *t
	for k, v := range data {
		tracerCopy.AddProperty(k, v)
	}
	t.logger.WithArguments(tracerCopy.getGroup()).Warn(msg)
}

func (t *tracer) Error(msg string) {
	t.logger.WithArguments(t.getGroup()).Error(msg)
}

func (t *tracer) Errorf(format string, data ...interface{}) {
	t.logger.WithArguments(t.getGroup()).Errorf(format, data...)
}

func (t *tracer) ErrorWithData(msg string, data map[string]any) {
	tracerCopy := *t
	for k, v := range data {
		tracerCopy.AddProperty(k, v)
	}
	t.logger.WithArguments(tracerCopy.getGroup()).Error(msg)
}

func (t *tracer) Debug(msg string) {
	t.logger.WithArguments(t.getGroup()).Debug(msg)
}

func (t *tracer) Debugf(format string, data ...interface{}) {
	t.logger.WithArguments(t.getGroup()).Debugf(format, data...)
}

func (t *tracer) DebugWithData(msg string, data map[string]any) {
	tracerCopy := *t
	for k, v := range data {
		tracerCopy.AddProperty(k, v)
	}
	t.logger.WithArguments(tracerCopy.getGroup()).Debug(msg)
}
