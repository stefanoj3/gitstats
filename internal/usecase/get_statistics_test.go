package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-github/github"
	"github.com/stefanoj3/gitstats/internal/usecase"
	"github.com/stretchr/testify/assert"
)

var baseTime = time.Date(2007, time.July, 7, 0, 0, 0, 0, time.UTC)

func TestGetStatisticsShouldCreateStatistics(t *testing.T) {
	pullRequestsFetcher := usecase.PullRequestsFetcherMock{
		FetchPullRequestsForFunc: func(
			ctx context.Context,
			from time.Time,
			to time.Time,
			organization string,
			repositories []string,
			usersHandles []string,
		) ([]*github.PullRequest, error) {
			prs := []*github.PullRequest{
				{
					CreatedAt: timeRef(baseTime.Add(-time.Hour)),
					MergedAt:  timeRef(baseTime),
				},
				{
					CreatedAt: timeRef(baseTime.Add(-time.Hour)),
					MergedAt:  timeRef(baseTime),
					ClosedAt:  timeRef(baseTime), // this should be ignored because it got merged
				},
				{
					CreatedAt: timeRef(baseTime.Add(-time.Hour)),
					ClosedAt:  timeRef(baseTime),
				},
				{
					CreatedAt: timeRef(baseTime.Add(-time.Hour)),
					ClosedAt:  timeRef(baseTime),
				},
				{
					CreatedAt: timeRef(baseTime.Add(-time.Hour)),
				},
			}

			return prs, nil
		},
	}
	commentsFetcher := usecase.CommentsFetcherMock{
		FetchCommentsForFunc: func(
			ctx context.Context,
			organization string,
			repository string,
			number int,
			usersHandles []string,
		) ([]*github.PullRequestComment, error) {
			return nil, nil
		},
	}

	sut := usecase.NewGetStatistics(&pullRequestsFetcher, &commentsFetcher)

	stats, err := sut.GetStatistics(
		context.Background(),
		baseTime.Add(-time.Hour*24*30),
		baseTime,
		"my-organization",
		[]string{"my-repository"},
		nil,
	)

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, stats.PullRequestsStatistics.Total, 5)
	assert.Equal(t, stats.PullRequestsStatistics.Merged, 2)
	assert.Equal(t, stats.PullRequestsStatistics.Closed, 2)
	assert.Equal(t, stats.PullRequestsStatistics.Open, 1)
}

func timeRef(t time.Time) *time.Time {
	return &t
}

func TestGetStatisticsShouldFailWhenPullRequestFetcherFails(t *testing.T) {
	errorMessage := "something 123"

	pullRequestsFetcher := usecase.PullRequestsFetcherMock{
		FetchPullRequestsForFunc: func(
			ctx context.Context,
			from time.Time,
			to time.Time,
			organization string,
			repositories []string,
			usersHandles []string,
		) ([]*github.PullRequest, error) {
			return nil, errors.New(errorMessage)
		},
	}
	commentsFetcher := usecase.CommentsFetcherMock{
		FetchCommentsForFunc: func(
			ctx context.Context,
			organization string,
			repository string,
			number int,
			usersHandles []string,
		) ([]*github.PullRequestComment, error) {
			return nil, nil
		},
	}

	sut := usecase.NewGetStatistics(&pullRequestsFetcher, &commentsFetcher)

	stats, err := sut.GetStatistics(
		context.Background(),
		baseTime.Add(-time.Hour*24*30),
		baseTime,
		"my-organization",
		[]string{"my-repository"},
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMessage)

	assert.Nil(t, stats)
}
