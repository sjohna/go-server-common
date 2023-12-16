package log

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultiplexLoggerWithSingleLogger(t *testing.T) {
	outBuffer := bytes.NewBuffer([]byte{})
	zLogger := zerolog.New(outBuffer).Level(zerolog.InfoLevel)

	t.Run("Single log", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.Info("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n", logged)
	})

	outBuffer.Reset()

	t.Run("Single log with field", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.WithField("key", "value").Info("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","key":"value","message":"test"}`+"\n", logged)
	})

	outBuffer.Reset()

	t.Run("Single log with fields", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.WithFields(Fields{"key1": "value1", "key2": "value2"}).Info("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","key1":"value1","key2":"value2","message":"test"}`+"\n", logged)
	})

	outBuffer.Reset()

	t.Run("Two logs, one below level", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.Info("test")
		logger.Debug("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n", logged)
	})

	outBuffer.Reset()

	t.Run("Array field", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.WithField("key", []string{"value1", "value2"}).Info("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","key":["value1","value2"],"message":"test"}`+"\n", logged)
	})

	outBuffer.Reset()

	t.Run("Struct field", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLogger})
		logger.WithField("key", struct {
			A string
			B string
		}{"value1", "value2"}).Info("test")
		logged := outBuffer.String()
		assert.Equal(t, `{"level":"info","key":{"A":"value1","B":"value2"},"message":"test"}`+"\n", logged)
	})
}

func TestMultiplexLoggerWithInfoAndDebugLoggers(t *testing.T) {
	outBufferInfo := bytes.NewBuffer([]byte{})
	zLoggerInfo := zerolog.New(outBufferInfo).Level(zerolog.InfoLevel)
	outBufferDebug := bytes.NewBuffer([]byte{})
	zLoggerDebug := zerolog.New(outBufferDebug).Level(zerolog.DebugLevel)

	t.Run("One info, one debug", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLoggerInfo, zLoggerDebug})
		logger.Info("test")
		logger.Debug("test")

		loggedInfo := outBufferInfo.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n", loggedInfo)

		loggedDebug := outBufferDebug.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n"+`{"level":"debug","message":"test"}`+"\n", loggedDebug)
	})

	outBufferInfo.Reset()
	outBufferDebug.Reset()

	t.Run("One info, one debug with field", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLoggerInfo, zLoggerDebug})
		logger.WithField("key1", "value1").Info("test")
		logger.WithField("key2", "value2").Debug("test")

		loggedInfo := outBufferInfo.String()
		assert.Equal(t, `{"level":"info","key1":"value1","message":"test"}`+"\n", loggedInfo)

		loggedDebug := outBufferDebug.String()
		assert.Equal(t, `{"level":"info","key1":"value1","message":"test"}`+"\n"+`{"level":"debug","key2":"value2","message":"test"}`+"\n", loggedDebug)
	})

	outBufferInfo.Reset()
	outBufferDebug.Reset()

	t.Run("Fatal and panic logs", func(t *testing.T) {
		logger := NewMultiplexLogger([]zerolog.Logger{zLoggerInfo, zLoggerDebug})
		logger.Fatal("test")
		logger.Panic("test")

		loggedInfo := outBufferInfo.String()
		assert.Equal(t, `{"level":"fatal","message":"test"}`+"\n"+`{"level":"panic","message":"test"}`+"\n", loggedInfo)

		loggedDebug := outBufferDebug.String()
		assert.Equal(t, `{"level":"fatal","message":"test"}`+"\n"+`{"level":"panic","message":"test"}`+"\n", loggedDebug)
	})
}

func TestCompoundLogger(t *testing.T) {
	outBufferInfo1 := bytes.NewBuffer([]byte{})
	zLoggerInfo1 := zerolog.New(outBufferInfo1).Level(zerolog.InfoLevel)
	outBufferDebug1 := bytes.NewBuffer([]byte{})
	zLoggerDebug1 := zerolog.New(outBufferDebug1).Level(zerolog.DebugLevel)
	logger1 := NewMultiplexLogger([]zerolog.Logger{zLoggerInfo1, zLoggerDebug1})

	outBufferInfo2 := bytes.NewBuffer([]byte{})
	zLoggerInfo2 := zerolog.New(outBufferInfo2).Level(zerolog.InfoLevel)
	logger2 := NewMultiplexLogger([]zerolog.Logger{zLoggerInfo2})

	logger := NewCompoundLogger([]Logger{logger1, logger2})

	t.Run("One info, one debug", func(t *testing.T) {
		logger.Info("test")
		logger.Debug("test")

		loggedInfo1 := outBufferInfo1.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n", loggedInfo1)

		loggedDebug1 := outBufferDebug1.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n"+`{"level":"debug","message":"test"}`+"\n", loggedDebug1)

		loggedInfo2 := outBufferInfo2.String()
		assert.Equal(t, `{"level":"info","message":"test"}`+"\n", loggedInfo2)
	})

	outBufferInfo1.Reset()
	outBufferDebug1.Reset()
	outBufferInfo2.Reset()

	t.Run("One info, one debug with field", func(t *testing.T) {
		logger.WithField("key1", "value1").Info("test")
		logger.WithField("key2", "value2").Debug("test")

		loggedInfo1 := outBufferInfo1.String()
		assert.Equal(t, `{"level":"info","key1":"value1","message":"test"}`+"\n", loggedInfo1)

		loggedDebug1 := outBufferDebug1.String()
		assert.Equal(t, `{"level":"info","key1":"value1","message":"test"}`+"\n"+`{"level":"debug","key2":"value2","message":"test"}`+"\n", loggedDebug1)

		loggedInfo2 := outBufferInfo2.String()
		assert.Equal(t, `{"level":"info","key1":"value1","message":"test"}`+"\n", loggedInfo2)
	})
}
