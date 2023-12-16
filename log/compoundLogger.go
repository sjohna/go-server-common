package log

type CompoundLogger struct {
	loggers []Logger
}

func NewCompoundLogger(loggers []Logger) CompoundLogger {
	return CompoundLogger{
		loggers,
	}
}

func (l CompoundLogger) WithField(key string, value interface{}) Logger {
	newLoggers := make([]Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.WithField(key, value)
	}
	return NewCompoundLogger(newLoggers)
}

func (l CompoundLogger) WithFields(fields map[string]interface{}) Logger {
	newLoggers := make([]Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.WithFields(fields)
	}
	return NewCompoundLogger(newLoggers)
}

func (l CompoundLogger) WithError(err error) Logger {
	newLoggers := make([]Logger, len(l.loggers))
	for i, logger := range l.loggers {
		newLoggers[i] = logger.WithError(err)
	}
	return NewCompoundLogger(newLoggers)
}

func (l CompoundLogger) Trace(msg string) {
	for _, logger := range l.loggers {
		logger.Trace(msg)
	}
}

func (l CompoundLogger) Tracef(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Tracef(format, v...)
	}
}

func (l CompoundLogger) Debug(msg string) {
	for _, logger := range l.loggers {
		logger.Debug(msg)
	}
}

func (l CompoundLogger) Debugf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Debugf(format, v...)
	}
}

func (l CompoundLogger) Info(msg string) {
	for _, logger := range l.loggers {
		logger.Info(msg)
	}
}

func (l CompoundLogger) Infof(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Infof(format, v...)
	}
}

func (l CompoundLogger) Warn(msg string) {
	for _, logger := range l.loggers {
		logger.Warn(msg)
	}
}

func (l CompoundLogger) Warnf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Warnf(format, v...)
	}
}

func (l CompoundLogger) Error(msg string) {
	for _, logger := range l.loggers {
		logger.Error(msg)
	}
}

func (l CompoundLogger) Errorf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Errorf(format, v...)
	}
}

func (l CompoundLogger) Panic(msg string) {
	for _, logger := range l.loggers {
		logger.Panic(msg)
	}
}

func (l CompoundLogger) Panicf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Panicf(format, v...)
	}
}

func (l CompoundLogger) Fatal(msg string) {
	for _, logger := range l.loggers {
		logger.Fatal(msg)
	}
}

func (l CompoundLogger) Fatalf(format string, v ...interface{}) {
	for _, logger := range l.loggers {
		logger.Fatalf(format, v...)
	}
}
