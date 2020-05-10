package usecase

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

type GithubDataFinder interface {
	FetchAllFor(
		ctx context.Context,
		from time.Time,
		to time.Time,
		delta time.Duration,
		organization string,
		repositories []string,
	) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error)
}
