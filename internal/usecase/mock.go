package usecase

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

type GithubDataFinderMock struct {
	FetchAllForFunc func(
		ctx context.Context,
		from time.Time,
		to time.Time,
		delta time.Duration,
		organization string,
		repositories []string,
	) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error)
}

func (g *GithubDataFinderMock) FetchAllFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	delta time.Duration,
	organization string,
	repositories []string,
) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error) {
	return g.FetchAllForFunc(ctx, from, to, delta, organization, repositories)
}
