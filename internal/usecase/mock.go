package usecase

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

type PullRequestsFetcherMock struct {
	FetchPullRequestsForFunc func(
		ctx context.Context,
		from time.Time,
		to time.Time,
		organization string,
		repositories []string,
		usersHandles []string,
	) ([]*github.PullRequest, error)
}

func (m *PullRequestsFetcherMock) FetchPullRequestsFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	organization string,
	repositories []string,
	usersHandles []string,
) ([]*github.PullRequest, error) {
	return m.FetchPullRequestsForFunc(ctx, from, to, organization, repositories, usersHandles)
}

type CommentsFetcherMock struct {
	FetchCommentsForFunc func(
		ctx context.Context,
		organization string,
		repository string,
		number int,
		usersHandles []string,
	) ([]*github.PullRequestComment, error)
}

func (m *CommentsFetcherMock) FetchCommentsFor(
	ctx context.Context,
	organization string,
	repository string,
	number int,
	usersHandles []string,
) ([]*github.PullRequestComment, error) {
	return m.FetchCommentsForFunc(ctx, organization, repository, number, usersHandles)
}
