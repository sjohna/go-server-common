package log

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
	Trace(msg string)
	Tracef(format string, v ...interface{})
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Error(msg string)
	Errorf(format string, v ...interface{})
	Fatal(msg string)
	Fatalf(format string, v ...interface{})
	Panic(msg string)
	Panicf(format string, v ...interface{})
}

type Fields map[string]interface{}
