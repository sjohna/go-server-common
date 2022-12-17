package service

import "github.com/sirupsen/logrus"

func ServiceFunctionLogger(log *logrus.Entry, serviceFunction string) *logrus.Entry {
	log = log.WithField("service-function", serviceFunction)
	log.Info("Service called")
	return log
}

func LogServiceReturn(log *logrus.Entry) {
	log.Info("Service returned")
}
