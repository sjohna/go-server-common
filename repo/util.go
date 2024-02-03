package repo

import (
	"context"
	"github.com/sjohna/go-server-common/log"
)

func RepoFunctionContext(ctx context.Context, repoFunction string) (repoContext context.Context, logger log.Logger) {
	logger = ctx.Value("logger").(log.Logger).WithField("repo-function", repoFunction)
	logger.Trace("Repo function called")
	repoContext = context.WithValue(ctx, "logger", logger)
	return
}

func LogRepoReturn(logger log.Logger) {
	logger.Trace("Repo function returned")
}
