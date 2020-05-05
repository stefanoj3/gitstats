package usecase

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/stefanoj3/gitstats/internal/domain/statistics"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func NewGetStatistics(githubDataFinder GithubDataFinder, logger *zap.Logger) *GetStatistics {
	return &GetStatistics{
		githubDataFinder: githubDataFinder,
		logger:           logger,
	}
}

type GetStatistics struct {
	githubDataFinder GithubDataFinder
	logger           *zap.Logger
}

func (g *GetStatistics) GetStatistics(
	ctx context.Context,
	from time.Time,
	to time.Time,
	delta time.Duration,
	organization string,
	repositories []string,
	userHandles []string,
) (*statistics.Statistics, error) {
	pullRequests, commits, comments, err := g.githubDataFinder.FetchAllFor(
		ctx,
		from,
		to,
		delta,
		organization,
		repositories,
	)
	if err != nil {
		return nil, errors.Wrap(err, "GetStatistics: failed to get github statistics")
	}

	stats := statistics.Statistics{
		PullRequestsStatistics: pullRequestStatisticsFromRawPullRequests(pullRequests, userHandles),
		UsersStatistics:        g.calculateUserStatistics(pullRequests, comments, commits, userHandles),
	}

	return &stats, nil
}

func (g *GetStatistics) calculateUserStatistics(
	pullRequests []*github.PullRequest,
	comments []*github.PullRequestComment,
	commits []*github.RepositoryCommit,
	userHandles []string,
) statistics.UsersStatistics {
	shouldExcludeFunc := buildUserFilterFunc(userHandles)

	stats := statistics.NewUserStatistics()

	for _, pr := range pullRequests {
		if shouldExcludeFunc(*pr.User.Login) {
			continue
		}

		stats.At(*pr.User.Login, *pr.CreatedAt).PullRequestsCreated++
	}

	for _, comment := range comments {
		if shouldExcludeFunc(*comment.User.Login) {
			continue
		}

		stats.At(*comment.User.Login, *comment.CreatedAt).Comments++
	}

	for _, commit := range commits {
		if shouldExcludeFunc(*commit.Author.Login) {
			continue
		}

		stats.At(*commit.Author.Login, *commit.Commit.Committer.Date).Commits++
	}

	return stats
}

func buildUserFilterFunc(userHandles []string) func(user string) bool {
	userMap := make(map[string]interface{}, len(userHandles))
	for _, u := range userHandles {
		userMap[u] = nil
	}

	shouldExclude := func(user string) bool {
		if len(userMap) == 0 {
			return false
		}

		_, ok := userMap[user]

		return !ok
	}

	return shouldExclude
}

func pullRequestStatisticsFromRawPullRequests(
	pullRequests []*github.PullRequest,
	userHandles []string,
) statistics.PullRequestsStatistics {
	shouldExcludeFunc := buildUserFilterFunc(userHandles)

	stats := statistics.PullRequestsStatistics{}

	timeToMergeDurations := make([]time.Duration, 0, stats.Total)

	for _, pr := range pullRequests {
		if shouldExcludeFunc(*pr.User.Login) {
			continue
		}

		stats.Total++

		if wasClosed(pr) {
			stats.Closed++
		}

		if warMerged(pr) {
			stats.Merged++

			timeToMerge := pr.MergedAt.Sub(*pr.CreatedAt)
			timeToMergeDurations = append(timeToMergeDurations, timeToMerge)
		}

		if isOpen(pr) {
			stats.Open++
		}
	}

	stats.TimeToMerge = averageDurations(timeToMergeDurations)

	return stats
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

func wasClosed(pullRequest *github.PullRequest) bool {
	return pullRequest.ClosedAt != nil && pullRequest.MergedAt == nil
}

func warMerged(pullRequest *github.PullRequest) bool {
	return pullRequest.MergedAt != nil
}

func isOpen(pullRequest *github.PullRequest) bool {
	return pullRequest.MergedAt == nil && pullRequest.ClosedAt == nil
}
