package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/google/go-github/github"
	"github.com/stefanoj3/gitstats/internal/usecase"
	"github.com/stretchr/testify/assert"
)

var (
	baseTime = time.Date(2007, time.July, 7, 0, 0, 0, 0, time.UTC)

	user1           = "user1"
	user2           = "user2"
	user3           = "user3"
	userToBeIgnored = "ignored_user_123"

	statisticsPRs = []*github.PullRequest{
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
			MergedAt:  timeRef(baseTime),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
			MergedAt:  timeRef(baseTime),
			ClosedAt:  timeRef(baseTime), // this should be ignored because it got merged
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
			ClosedAt:  timeRef(baseTime),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
			ClosedAt:  timeRef(baseTime),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime),
		},
		{
			User:      &github.User{Login: &userToBeIgnored},
			CreatedAt: timeRef(baseTime),
		},
	}
	statisticsCommits = []*github.RepositoryCommit{
		{
			Author: &github.User{Login: &user1},
			Commit: &github.Commit{Committer: &github.CommitAuthor{Date: timeRef(baseTime)}},
		},
		{
			Author: &github.User{Login: &user1},
			Commit: &github.Commit{Committer: &github.CommitAuthor{Date: timeRef(baseTime.Add(time.Hour))}},
		},
		{
			Author: &github.User{Login: &user2},
			Commit: &github.Commit{Committer: &github.CommitAuthor{Date: timeRef(baseTime)}},
		},
		{
			Author: &github.User{Login: &userToBeIgnored},
			Commit: &github.Commit{Committer: &github.CommitAuthor{Date: timeRef(baseTime)}},
		},
	}
	statisticsComments = []*github.PullRequestComment{
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour * 5)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour * 12)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour * 3)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour * 7)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(-time.Hour)),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime),
		},
		{
			User:      &github.User{Login: &user1},
			CreatedAt: timeRef(baseTime.Add(+time.Hour)),
		},
		{
			User:      &github.User{Login: &userToBeIgnored},
			CreatedAt: timeRef(baseTime.Add(+time.Hour)),
		},
	}
)

func TestGetStatisticsShouldCreateStatistics(t *testing.T) {
	finder := usecase.GithubDataFinderMock{FetchAllForFunc: func(
		ctx context.Context,
		from time.Time,
		to time.Time,
		delta time.Duration,
		organization string,
		repositories []string,
	) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error) {
		return statisticsPRs, statisticsCommits, statisticsComments, nil
	}}

	sut := usecase.NewGetStatistics(&finder, zap.NewNop())

	stats, err := sut.GetStatistics(
		context.Background(),
		baseTime.Add(-time.Hour*24*30),
		baseTime,
		time.Microsecond,
		"my-organization",
		[]string{"my-repository"},
		[]string{user1, user2, user3},
	)

	assert.NoError(t, err)
	assert.NotNil(t, stats)

	assert.Equal(t, 6, stats.PullRequestsStatistics.Total)
	assert.Equal(t, 2, stats.PullRequestsStatistics.Merged)
	assert.Equal(t, 2, stats.PullRequestsStatistics.Closed)
	assert.Equal(t, 2, stats.PullRequestsStatistics.Open)

	assert.Equal(t, 1, stats.UsersStatistics.At(user1, baseTime).PullRequestsCreated)
	assert.Equal(t, 5, stats.UsersStatistics.At(user1, baseTime.Add(-time.Hour)).PullRequestsCreated)

	assert.Equal(t, 0, stats.UsersStatistics.At(user2, baseTime).PullRequestsCreated)
	assert.Equal(t, 0, stats.UsersStatistics.At(user2, baseTime.Add(-time.Hour)).PullRequestsCreated)

	assert.Equal(t, 0, stats.UsersStatistics.At(user3, baseTime).PullRequestsCreated)
	assert.Equal(t, 0, stats.UsersStatistics.At(user3, baseTime.Add(-time.Hour)).PullRequestsCreated)

	assert.Equal(t, 2, stats.UsersStatistics.At(user1, baseTime).Commits)
	assert.Equal(t, 1, stats.UsersStatistics.At(user2, baseTime).Commits)
	assert.Equal(t, 0, stats.UsersStatistics.At(user3, baseTime).Commits)

	assert.Equal(t, 5, stats.UsersStatistics.At(user1, baseTime.Add(-time.Hour)).Comments)
	assert.Equal(t, 2, stats.UsersStatistics.At(user1, baseTime).Comments)

	assert.Equal(t, 0, stats.UsersStatistics.At(user2, baseTime).Comments)
	assert.Equal(t, 0, stats.UsersStatistics.At(user3, baseTime).Comments)
}

func TestGetStatisticsShouldFailWhenPullRequestFetcherFails(t *testing.T) {
	errorMessage := "something 123"

	finder := usecase.GithubDataFinderMock{FetchAllForFunc: func(
		ctx context.Context,
		from time.Time,
		to time.Time,
		delta time.Duration,
		organization string,
		repositories []string,
	) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error) {
		return nil, nil, nil, errors.New(errorMessage)
	}}

	sut := usecase.NewGetStatistics(&finder, zap.NewNop())

	stats, err := sut.GetStatistics(
		context.Background(),
		baseTime.Add(-time.Hour*24*30),
		baseTime,
		time.Microsecond,
		"my-organization",
		[]string{"my-repository"},
		nil,
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errorMessage)

	assert.Nil(t, stats)
}

func timeRef(t time.Time) *time.Time {
	return &t
}
