package service

import (
	"context"
	"github.com/sjohna/go-server-common/log"
)

func ServiceFunctionContext(ctx context.Context, serviceFunction string) (serviceContext context.Context, logger log.Logger) {
	logger = ctx.Value("logger").(log.Logger).WithField("service-function", serviceFunction)
	logger.Debug("Service function called")
	serviceContext = context.WithValue(ctx, "logger", logger)
	return
}

func LogServiceReturn(logger log.Logger) {
	logger.Trace("Service function returned")
}
