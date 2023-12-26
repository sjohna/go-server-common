package main

import "github.com/sjohna/go-server-common/log"

func main() {
	logger, configLogger := log.GetApplicationLoggers("testOutput/serverCommonTest", "go-server-common-test")

	logger.Info("Hello world")
	configLogger.Info("Hello world")

	logger.WithField("test", "test").Panic("Test primitive field")
	logger.WithField("test", struct{ Test string }{"test"}).Panic("Test struct field")
	logger.WithField("test", map[string]string{"test": "test"}).Panic("Test map field")
	logger.WithField("test", []string{"test", "test"}).Panic("Test slice field")
}
