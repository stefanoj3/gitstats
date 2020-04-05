package usecase

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

type PullRequestsFetcher interface {
	FetchPullRequestsFor(
		ctx context.Context,
		from time.Time,
		to time.Time,
		organization string,
		repositories []string,
		usersHandles []string,
	) ([]*github.PullRequest, error)
}

type CommentsFetcher interface {
	FetchCommentsFor(
		ctx context.Context,
		organization string,
		repository string,
		number int,
		usersHandles []string,
	) ([]*github.PullRequestComment, error)
}
