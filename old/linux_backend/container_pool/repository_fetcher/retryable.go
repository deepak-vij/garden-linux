package repository_fetcher

import (
	"github.com/cloudfoundry-incubator/garden-linux/process"
	"github.com/pivotal-golang/lager"
)

type Retryable struct {
	RepositoryFetcher
}

func (retryable Retryable) Fetch(logger lager.Logger, repoName string, tag string) (string, process.Env, []string, error) {
	var res string
	var err error
	var envvars process.Env
	var volumes []string

	for attempt := 1; attempt <= 3; attempt++ {
		res, envvars, volumes, err = retryable.RepositoryFetcher.Fetch(logger, repoName, tag)
		if err == nil {
			break
		}

		logger.Error("failed-to-fetch", err, lager.Data{
			"attempt": attempt,
			"of":      3,
		})
	}

	return res, envvars, volumes, err
}
