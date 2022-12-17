package repo

func RepoFunctionLogger(log *logrus.Entry, repoFunction string) *logrus.Entry {
	log = log.WithField("repo-function", repoFunction)
	log.Info("Repo called")
	return log
}

func LogRepoReturn(log *logrus.Entry) {
	log.Info("Repo returned")
}
