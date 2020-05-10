package git

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const pullRequestsPerPage = 100

func NewGithubPullRequestsRepository(client *github.Client, logger *zap.Logger) *GithubPullRequestsRepository {
	return &GithubPullRequestsRepository{client: client, logger: logger}
}

type GithubPullRequestsRepository struct {
	client *github.Client
	logger *zap.Logger
}

func (r *GithubPullRequestsRepository) FindPullRequestsFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	organization string,
	repositories []string,
) ([]*github.PullRequest, error) {
	result := make([]*github.PullRequest, 0)

	for _, repository := range repositories {
		prs, err := r.fetchAllFor(ctx, from, to, organization, repository)
		if err != nil {
			return nil, err
		}

		result = append(result, prs...)
	}

	return result, nil
}

func (r *GithubPullRequestsRepository) fetchAllFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	organization string,
	repository string,
) ([]*github.PullRequest, error) {
	listOptions := github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: pullRequestsPerPage,
		},
		State:     "all",
		Sort:      "created",
		Direction: "desc",
	}

	result := make([]*github.PullRequest, 0, pullRequestsPerPage)

	shouldRun := true
	for shouldRun {
		r.logger.Debug(
			"listing commits",
			zap.String("organization", organization),
			zap.String("repository", organization),
			zap.Int("page", listOptions.Page),
			zap.Int("perPage", listOptions.PerPage),
		)

		pullRequests, _, err := r.client.PullRequests.List(ctx, organization, repository, &listOptions)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pull requests")
		}

		if len(pullRequests) == 0 {
			break
		}

		for _, pr := range pullRequests {
			if pr.CreatedAt.Before(from) {
				// we fetch pull requests by created_at desc, as soon as we find one older
				// than our from(time) we can stop executing calls
				shouldRun = false
				continue
			}

			if pr.CreatedAt.After(to) {
				continue
			}

			result = append(result, pr)
		}

		listOptions.Page++
	}

	return result, nil
}
