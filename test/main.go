package main

import (
	"context"
	"github.com/sjohna/go-server-common/handler"
	"github.com/sjohna/go-server-common/log"
	"github.com/sjohna/go-server-common/repo"
	"github.com/sjohna/go-server-common/service"
	"net/http"
)

func main() {
	logger, configLogger := log.GetApplicationLoggers("testOutput/serverCommonTest", "go-server-common-test")

	logger.Info("Hello world")
	configLogger.Info("Hello world")

	logger.WithField("test", "test").Panic("Test primitive field")
	logger.WithField("test", struct{ Test string }{"test"}).Panic("Test struct field")
	logger.WithField("test", map[string]string{"test": "test"}).Panic("Test map field")
	logger.WithField("test", []string{"test", "test"}).Panic("Test slice field")

	testContext := context.WithValue(context.Background(), "logger", logger)

	testRequest, _ := http.NewRequestWithContext(testContext, "GET", "http://localhost:8080", nil)

	TestHandlerFunction(testRequest)
}

func TestHandlerFunction(req *http.Request) {
	handlerContext, logger := handler.HandlerContext(req, "TestHandlerFunction")
	defer handler.LogHandlerReturn(logger)

	logger.Info("In TestHandlerFunction")

	TestServiceFunction(handlerContext)
}

func TestServiceFunction(ctx context.Context) {
	serviceContext, logger := service.ServiceFunctionContext(ctx, "TestServiceFunction")
	defer service.LogServiceReturn(logger)

	logger.Info("In TestServiceFunction")

	TestRepoFunction(serviceContext)
}

func TestRepoFunction(ctx context.Context) {
	_, logger := repo.RepoFunctionContext(ctx, "TestRepoFunction")
	defer repo.LogRepoReturn(logger)

	logger.Info("In TestRepoFunction")
}
