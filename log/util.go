package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"path"
	"time"
)
import "gopkg.in/natefinch/lumberjack.v2"

func GetApplicationLoggers(logDirectory string, applicationName string) (logger Logger, configLogger Logger) {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger = NewMultiplexLogger([]zerolog.Logger{
		newZeroLogger(path.Join(logDirectory, applicationName+"_ERROR.log"), zerolog.ErrorLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_WARN.log"), zerolog.WarnLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_INFO.log"), zerolog.InfoLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_DEBUG.log"), zerolog.DebugLevel).With().Stack().Logger(),
		newZeroLoggerWithStandardOut(path.Join(logDirectory, applicationName+"_TRACE.log"), zerolog.TraceLevel).With().Stack().Logger(),
	})

	configBaseLogger := NewMultiplexLogger([]zerolog.Logger{
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_ERROR.log"), zerolog.ErrorLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_WARN.log"), zerolog.WarnLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_INFO.log"), zerolog.InfoLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_DEBUG.log"), zerolog.DebugLevel).With().Stack().Logger(),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_TRACE.log"), zerolog.TraceLevel).With().Stack().Logger(),
	})

	configLogger = NewCompoundLogger([]Logger{configBaseLogger, logger})

	return
}

func newZeroLogger(filePath string, level zerolog.Level) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     36500, // 100 years
		Compress:   false,
	}

	return zerolog.New(fileLogger).Level(level).With().Timestamp().Logger()
}

func newZeroLoggerWithStandardOut(filePath string, level zerolog.Level) zerolog.Logger {
	fileLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     36500, // 100 years
		Compress:   false,
	}

	io.MultiWriter(fileLogger, os.Stdout)

	return zerolog.New(io.MultiWriter(fileLogger, os.Stdout)).Level(level).With().Timestamp().Logger()
}
