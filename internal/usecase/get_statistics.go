package usecase

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type UserStatistics struct {
	// nolint:godox
	// TODO: define and add statistics by user
}

type PullRequestsStatistics struct {
	// TimeToMerge represents the average time it takes for a PR to be merged
	TimeToMerge time.Duration
	// Merged represents the total or PRs that got merged
	Merged int
	// Closed represents the total or PRs that got closed but never merged
	Closed int
	// Open represents the total or PRs that are still open
	Open int
	// Total is the total amount of PRs
	Total int
}

type Statistics struct {
	PullRequestsStatistics PullRequestsStatistics
	UsersStatistics        []UserStatistics
}

func NewGetStatistics(pullRequestsFetcher PullRequestsFetcher, commentsFetcher CommentsFetcher) *GetStatistics {
	return &GetStatistics{
		pullRequestsFetcher: pullRequestsFetcher,
		commentsFetcher:     commentsFetcher,
	}
}

type GetStatistics struct {
	pullRequestsFetcher PullRequestsFetcher
	commentsFetcher     CommentsFetcher
}

func (g *GetStatistics) GetStatistics(
	ctx context.Context,
	from time.Time,
	to time.Time,
	organization string,
	repositories []string,
	usersHandles []string,
) (*Statistics, error) {
	pullRequests, err := g.pullRequestsFetcher.FetchPullRequestsFor(
		ctx,
		from,
		to,
		organization,
		repositories,
		usersHandles,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetStatistics: failed to Get statistics")
	}

	statistics := Statistics{
		PullRequestsStatistics: pullRequestStatisticsFromRawPullRequests(pullRequests),
	}

	return &statistics, nil
}

func pullRequestStatisticsFromRawPullRequests(pullRequests []*github.PullRequest) PullRequestsStatistics {
	statistics := PullRequestsStatistics{}

	statistics.Total = len(pullRequests)

	timeToMergeDurations := make([]time.Duration, 0, statistics.Total)

	for _, pr := range pullRequests {
		if pr.ClosedAt != nil && pr.MergedAt == nil {
			statistics.Closed++
		}

		if pr.MergedAt != nil {
			statistics.Merged++

			timeToMerge := pr.MergedAt.Sub(*pr.CreatedAt)
			timeToMergeDurations = append(timeToMergeDurations, timeToMerge)
		}

		if pr.MergedAt == nil && pr.ClosedAt == nil {
			statistics.Open++
		}
	}

	statistics.TimeToMerge = averageDurations(timeToMergeDurations)

	return statistics
}

func averageDurations(values []time.Duration) time.Duration {
	count := len(values)
	if count == 0 {
		return 0
	}

	var total time.Duration
	for _, value := range values {
		total += value
	}

	return total / time.Duration(count)
}
