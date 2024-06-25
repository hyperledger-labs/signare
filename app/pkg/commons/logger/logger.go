// Package logger defines the logging utilities for the signare.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	signertime "github.com/hyperledger-labs/signare/app/pkg/commons/time"
	"github.com/hyperledger-labs/signare/app/pkg/entities"
)

// Logger defines the functionality to log with different levels in the application.
type Logger interface {
	// WithArguments injects extra arguments to the log entry
	WithArguments(args ...any) Logger
	// Info creates a log entry with LevelInfo
	Info(msg string)
	// Infof creates a formatted log entry with LevelInfo
	Infof(format string, args ...any)
	// Debug creates a log entry with LevelDebug
	Debug(msg string)
	// Debugf creates a formatted log entry with LevelDebug
	Debugf(format string, args ...any)
	// Warn creates a log entry with LevelWarn
	Warn(msg string)
	// Warnf creates a formatted log entry with LevelWarn
	Warnf(format string, args ...any)
	// Error creates a log entry with LevelError
	Error(msg string)
	// Errorf creates a formatted log entry with LevelError
	Errorf(format string, args ...interface{})
}

type logger struct {
	*slog.Logger
	ctx             context.Context
	arguments       []any
	ctxKeysRegistry map[entities.ContextKey]LogKey
}

type LogKey string

var (
	_             Logger = new(logger)
	defaultLogger *logger
)

func init() {
	levelInfo := LevelInfo
	RegisterLogger(Options{
		Level:     &levelInfo,
		LogOutput: os.Stdout,
	})
}

// Options configures a Logger.
type Options struct {
	// Level the log level
	Level *Level
	// LogOutput where the logs are written to
	LogOutput io.Writer
	// CtxKeysRegistry defines which keys from the context, if present, will be logged as a default attribute.
	CtxKeysRegistry map[entities.ContextKey]LogKey
}

// RegisterLogger for the application with the provided options.
func RegisterLogger(options Options) {
	slogLevel := slog.LevelInfo
	if options.Level != nil {
		if l, ok := logLevelTranslator[*options.Level]; ok {
			slogLevel = l
		}
	}

	slogOptions := &slog.HandlerOptions{
		Level:       slogLevel,
		ReplaceAttr: replaceAttributes,
	}
	var logHandler slog.Handler = slog.NewJSONHandler(options.LogOutput, slogOptions)
	if isLocalMode() {
		logHandler = slog.NewTextHandler(options.LogOutput, slogOptions)
	}
	defaultLogger = &logger{
		Logger:          slog.New(logHandler),
		ctxKeysRegistry: options.CtxKeysRegistry,
	}
}

// ToLevel maps a string to its corresponding Level returning false if it's not defined.
func ToLevel(level string) (*Level, bool) {
	for l := range logLevelTranslator {
		if string(l) == strings.ToUpper(level) {
			return &l, true
		}
	}
	return nil, false
}

// replaceAttributes changes the default keys of the slog built-in attributes.
func replaceAttributes(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		if isLocalMode() {
			return slog.Attr{
				Key:   slog.TimeKey,
				Value: slog.StringValue(time.Now().Format(time.DateTime)),
			}
		}
		return slog.Attr{
			Key:   slog.TimeKey,
			Value: slog.Int64Value(signertime.Now().ToInt64()),
		}
	case slog.MessageKey:
		return slog.Attr{
			Key:   "message",
			Value: a.Value,
		}
	}
	return slog.Attr{Key: a.Key, Value: a.Value}
}

// LogEntry returns a new Logger
func LogEntry(ctx context.Context) Logger {
	l := &logger{
		Logger:          defaultLogger.Logger,
		ctx:             ctx,
		ctxKeysRegistry: defaultLogger.ctxKeysRegistry,
	}
	return l.injectContextKeys()
}

func (l *logger) injectContextKeys() *logger {
	for ctxKey, logKey := range l.ctxKeysRegistry {
		value := l.ctx.Value(ctxKey)
		if value == nil {
			continue
		}
		valueStr := fmt.Sprintf("%v", value)
		l.arguments = append(l.arguments, slog.String(string(logKey), valueStr))
	}
	return l
}

func (l *logger) WithArguments(args ...any) Logger {
	l.arguments = args
	l.injectContextKeys()
	return l
}

func (l *logger) Info(msg string) {
	l.Logger.With(l.arguments...).InfoContext(l.ctx, msg)
}

func (l *logger) Infof(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	l.Logger.With(l.arguments...).InfoContext(l.ctx, m)
}

func (l *logger) Debug(msg string) {
	l.Logger.With(l.arguments...).DebugContext(l.ctx, msg)
}

func (l *logger) Debugf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	l.Logger.With(l.arguments...).DebugContext(l.ctx, m)
}

func (l *logger) Warn(msg string) {
	l.Logger.With(l.arguments...).WarnContext(l.ctx, msg)
}

func (l *logger) Warnf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	l.Logger.With(l.arguments...).WarnContext(l.ctx, m)
}

func (l *logger) Error(msg string) {
	l.Logger.With(l.arguments...).ErrorContext(l.ctx, msg)
}

func (l *logger) Errorf(format string, args ...any) {
	m := fmt.Sprintf(format, args...)
	l.Logger.With(l.arguments...).ErrorContext(l.ctx, m)
}
