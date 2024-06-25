package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"
	"github.com/hyperledger-labs/signare/app/pkg/entities"

	"github.com/stretchr/testify/assert"
)

var (
	msg        = "this is important"
	levelDebug = logger.LevelDebug
	levelInfo  = logger.LevelInfo

	ctx        = context.Background()
	contextKey = entities.ContextKey("Context-Key")
	buffer     = &bytes.Buffer{}
)

const important = "important"

func resetLogger() {
	buffer.Reset()
	logger.RegisterLogger(logger.Options{
		Level:     &levelDebug,
		LogOutput: buffer,
		CtxKeysRegistry: map[entities.ContextKey]logger.LogKey{
			contextKey: "contextKey",
		},
	})
}

func TestProvideDefaultLogger(t *testing.T) {
	t.Setenv("DEV_LOCAL_MODE", "no")
	type logEntry struct {
		Time    int64        `json:"time"`
		Level   logger.Level `json:"level"`
		Message string       `json:"message"`
	}

	t.Run("JSON Handler is enabled by default", func(t *testing.T) {
		resetLogger()

		logMessage := "This is the log message"
		logger.LogEntry(ctx).Info(logMessage)
		logOutput := logEntry{}
		err := json.Unmarshal(buffer.Bytes(), &logOutput)
		assert.NoError(t, err)
		assert.Equal(t, logMessage, logOutput.Message)
		assert.Equal(t, logger.LevelInfo, logOutput.Level)
		assert.NotZero(t, logOutput.Time)
	})

	t.Run("Text Handler is enabled in local mode", func(t *testing.T) {
		defer resetLogger()
		t.Setenv("DEV_LOCAL_MODE", "yes")
		levelDebug := logger.LevelDebug
		logger.RegisterLogger(logger.Options{
			Level:     &levelDebug,
			LogOutput: buffer,
		})

		logMessage := "This is the log message"
		logger.LogEntry(ctx).Info(logMessage)
		logOutput := logEntry{}

		expectedOutput := `level=INFO message="This is the log message"`
		err := json.Unmarshal(buffer.Bytes(), &logOutput)
		assert.Error(t, err)
		assert.Contains(t, buffer.String(), expectedOutput)
	})

	t.Run("Lower levels than the configured one are not printed", func(t *testing.T) {
		defer resetLogger()
		logger.RegisterLogger(logger.Options{
			Level:     &levelInfo,
			LogOutput: buffer,
		})

		logMessage := "This is a debug message"
		logger.LogEntry(ctx).Debug(logMessage)
		logOutput := logEntry{}
		err := json.Unmarshal(buffer.Bytes(), &logOutput)
		assert.Error(t, err)
	})
}

func TestLogger_withoutArguments(t *testing.T) {
	t.Setenv("DEV_LOCAL_MODE", "no")
	resetLogger()

	t.Run("INFO Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"INFO","message":"this is important"}`

		logger.LogEntry(ctx).Info(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"INFO","message":"this is an important message"}`

		logger.LogEntry(ctx).Infof("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("WARN Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"WARN","message":"this is important"}`

		logger.LogEntry(ctx).Warn(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"WARN","message":"this is an important message"}`

		logger.LogEntry(ctx).Warnf("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("DEBUG Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"DEBUG","message":"this is important"}`

		logger.LogEntry(ctx).Debug(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"DEBUG","message":"this is an important message"}`

		logger.LogEntry(ctx).Debugf("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("ERROR Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"ERROR","message":"this is important"}`

		logger.LogEntry(ctx).Error(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})
}

func TestLogger_withArguments(t *testing.T) {
	t.Setenv("DEV_LOCAL_MODE", "no")
	resetLogger()

	t.Run("INFO Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"INFO","message":"this is important","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Info(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"INFO","message":"this is an important message","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Infof("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("WARN Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"WARN","message":"this is important","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Warn(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"WARN","message":"this is an important message","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Warnf("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("DEBUG Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"DEBUG","message":"this is important","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Debug(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"DEBUG","message":"this is an important message","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Debugf("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("ERROR Level", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"ERROR","message":"this is important","threadID":10}`

		logger.LogEntry(ctx).WithArguments("threadID", 10).Error(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})
}

func TestLogger_withContextKeyValues(t *testing.T) {
	t.Setenv("DEV_LOCAL_MODE", "no")
	resetLogger()
	ctxWithValue := context.WithValue(ctx, contextKey, "test-value")

	t.Run("Sucess with arguments and context values", func(t *testing.T) {
		defer resetLogger()
		expected := `{"level":"INFO","message":"this is important","threadID":10,"contextKey":"test-value"}`

		logger.LogEntry(ctxWithValue).WithArguments("threadID", 10).Info(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		expected = `{"level":"INFO","message":"this is an important message","threadID":10,"contextKey":"test-value"}`

		logger.LogEntry(ctxWithValue).WithArguments("threadID", 10).Infof("this is an %s message", important)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})
}

func logWithoutTimeAttr(b *bytes.Buffer) string {
	log := b.String()
	timestampPrefixLength := len(`{"time":1709634409179,`)
	log = fmt.Sprintf("{%s", log[timestampPrefixLength:])
	return strings.TrimSuffix(log, "\n")
}
