package logger_test

import (
	"testing"
	"time"

	"github.com/hyperledger-labs/signare/app/pkg/commons/logger"

	"github.com/stretchr/testify/assert"
)

func TestTracer(t *testing.T) {
	t.Setenv("DEV_LOCAL_MODE", "no")
	resetLogger()

	tracer := logger.NewTracer(ctx)
	tracer.AddProperty("testString", "lorem ipsum")
	tracer.AddProperty("testInt", -50)
	tracer.AddProperty("testInt64", int64(-30))
	tracer.AddProperty("testUint64", uint64(100))
	tracer.AddProperty("testFloat64", 12.50)
	tracer.AddProperty("testBool", true)
	tracer.AddProperty("testTimeDuration", time.Millisecond)

	t.Run("INFO Level", func(t *testing.T) {
		// this properties must always be present
		expected := `{"level":"INFO","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Trace(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// "data" must be a temporary tracer's object extra argument
		expected = `{"level":"INFO","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000,"data":"extraData"}}`
		tracer.TraceWithData(msg, map[string]any{
			"data": "extraData",
		})
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// only tracer's permanent properties must be present
		expected = `{"level":"INFO","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Trace(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("WARN Level", func(t *testing.T) {
		// this properties must always be present
		expected := `{"level":"WARN","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Warn(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// "data" must be a temporary tracer's object extra argument
		expected = `{"level":"WARN","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000,"data":"extraData"}}`
		tracer.WarnWithData(msg, map[string]any{
			"data": "extraData",
		})
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// only tracer's permanent properties must be present
		expected = `{"level":"WARN","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Warn(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("DEBUG Level", func(t *testing.T) {
		// this properties must always be present
		expected := `{"level":"DEBUG","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Debug(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// "data" must be a temporary tracer's object extra argument
		expected = `{"level":"DEBUG","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000,"data":"extraData"}}`
		tracer.DebugWithData(msg, map[string]any{
			"data": "extraData",
		})
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// only tracer's permanent properties must be present
		expected = `{"level":"DEBUG","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Debug(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})

	t.Run("ERROR Level", func(t *testing.T) {
		// this properties must always be present
		expected := `{"level":"ERROR","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Error(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// "data" must be a temporary tracer's object extra argument
		expected = `{"level":"ERROR","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000,"data":"extraData"}}`
		tracer.ErrorWithData(msg, map[string]any{
			"data": "extraData",
		})
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()

		// only tracer's permanent properties must be present
		expected = `{"level":"ERROR","message":"this is important","properties":{"testString":"lorem ipsum","testInt":-50,"testInt64":-30,"testUint64":100,"testFloat64":12.5,"testBool":true,"testTimeDuration":1000000}}`
		tracer.Error(msg)
		assert.Equal(t, expected, logWithoutTimeAttr(buffer))
		buffer.Reset()
	})
}
