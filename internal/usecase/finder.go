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

type PullRequestsFinder interface {
	FindPullRequestsFor(
		ctx context.Context,
		from time.Time,
		to time.Time,
		organization string,
		repositories []string,
	) ([]*github.PullRequest, error)
}

type CommentsFinder interface {
	FindCommentsFor(
		ctx context.Context,
		organization string,
		repository string,
		number int,
	) ([]*github.PullRequestComment, error)
}

type CommitsFinder interface {
	FindCommitsFor(
		ctx context.Context,
		organization string,
		repository string,
		number int,
	) ([]*github.RepositoryCommit, error)
}
