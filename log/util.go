package log

import (
	"github.com/rs/zerolog"
	"io"
	"os"
	"path"
)
import "gopkg.in/natefinch/lumberjack.v2"

func GetApplicationLoggers(logDirectory string, applicationName string) (logger Logger, configLogger Logger) {
	logger = NewMultiplexLogger([]zerolog.Logger{
		newZeroLogger(path.Join(logDirectory, applicationName+"_ERROR.log"), zerolog.ErrorLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_WARN.log"), zerolog.WarnLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_INFO.log"), zerolog.InfoLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_DEBUG.log"), zerolog.DebugLevel),
		newZeroLoggerWithStandardOut(path.Join(logDirectory, applicationName+"_TRACE.log"), zerolog.TraceLevel),
	})

	configBaseLogger := NewMultiplexLogger([]zerolog.Logger{
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_ERROR.log"), zerolog.ErrorLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_WARN.log"), zerolog.WarnLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_INFO.log"), zerolog.InfoLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_DEBUG.log"), zerolog.DebugLevel),
		newZeroLogger(path.Join(logDirectory, applicationName+"_CONFIG_TRACE.log"), zerolog.TraceLevel),
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
