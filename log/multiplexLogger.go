package log

import (
	"github.com/rs/zerolog"
	"github.com/sjohna/go-server-common/errors"
)

type MultiplexLogger struct {
	loggers []zerolog.Logger
}

func NewMultiplexLogger(loggers []zerolog.Logger) MultiplexLogger {
	return MultiplexLogger{
		loggers,
	}
}

func (l MultiplexLogger) WithField(key string, value interface{}) Logger {
	newLoggers := make([]zerolog.Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.With().Interface(key, value).Logger()
	}
	return NewMultiplexLogger(newLoggers)
}

func (l MultiplexLogger) WithFields(fields map[string]interface{}) Logger {
	newLoggers := make([]zerolog.Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.With().Fields(fields).Logger()
	}
	return NewMultiplexLogger(newLoggers)
}

func (l MultiplexLogger) WithError(err errors.Error) Logger {
	newLoggers := make([]zerolog.Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.With().Err(err).Logger()

		if queryErr, isQueryErr := err.(*errors.QueryError); isQueryErr {
			newLoggers[i] = logger.With().Fields(map[string]interface{}{
				"origin":     errors.OriginString(queryErr.Origin),
				"errorStack": queryErr.StackTrace,
				"innerError": queryErr.Inner,
				"query":      queryErr.Query,
				"queryArgs":  queryErr.Args,
			}).Logger()
		} else if appErr, isAppErr := err.(*errors.ApplicationError); isAppErr {
			newLoggers[i] = logger.With().Fields(map[string]interface{}{
				"origin":     errors.OriginString(appErr.Origin),
				"errorStack": appErr.StackTrace,
				"innerError": appErr.Inner,
			}).Logger()
		}
	}
	return NewMultiplexLogger(newLoggers)
}

func (l MultiplexLogger) Trace(msg string) {
	for _, logger := range l.loggers {
		logger.Trace().Msg(msg)
	}
}

func (l MultiplexLogger) Tracef(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Trace().Msgf(format, v...)
	}
}

func (l MultiplexLogger) Debug(msg string) {
	for _, logger := range l.loggers {
		logger.Debug().Msg(msg)
	}
}

func (l MultiplexLogger) Debugf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Debug().Msgf(format, v...)
	}
}

func (l MultiplexLogger) Info(msg string) {
	for _, logger := range l.loggers {
		logger.Info().Msg(msg)
	}
}

func (l MultiplexLogger) Infof(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Info().Msgf(format, v...)
	}
}

func (l MultiplexLogger) Warn(msg string) {
	for _, logger := range l.loggers {
		logger.Warn().Msg(msg)
	}
}

func (l MultiplexLogger) Warnf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Warn().Msgf(format, v...)
	}
}

func (l MultiplexLogger) Error(msg string) {
	for _, logger := range l.loggers {
		logger.Error().Msg(msg)
	}
}

func (l MultiplexLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Error().Msgf(format, v...)
	}
}

func (l MultiplexLogger) Panic(msg string) {
	for _, logger := range l.loggers {
		logger.WithLevel(zerolog.PanicLevel).Msg(msg) // doing it this way so that this doesn't actually kill the goroutine
	}
}

func (l MultiplexLogger) Panicf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.WithLevel(zerolog.PanicLevel).Msgf(format, v...) // doing it this way so that this doesn't actually kill the goroutine
	}
}

func (l MultiplexLogger) Fatal(msg string) {
	for _, logger := range l.loggers {
		logger.WithLevel(zerolog.FatalLevel).Msg(msg) // doing it this way so that this doesn't actually kill the process
	}
}

func (l MultiplexLogger) Fatalf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.WithLevel(zerolog.FatalLevel).Msgf(format, v...) // doing it this way so that this doesn't actually kill the process
	}
}
