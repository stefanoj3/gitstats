package git

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const commitsPerPage = 100

func NewGithubCommitsRepository(client *github.Client, logger *zap.Logger) *GithubCommitsRepository {
	return &GithubCommitsRepository{client: client, logger: logger}
}

type GithubCommitsRepository struct {
	client *github.Client
	logger *zap.Logger
}

func (r *GithubCommitsRepository) FindCommitsFor(
	ctx context.Context,
	organization string,
	repository string,
	number int,
) ([]*github.RepositoryCommit, error) {
	opts := github.ListOptions{PerPage: commitsPerPage, Page: 1}

	result := make([]*github.RepositoryCommit, 0, commitsPerPage)

	shouldRun := true

	//nolint:dupl
	for shouldRun {
		r.logger.Debug(
			"listing commits",
			zap.String("organization", organization),
			zap.String("repository", organization),
			zap.Int("number", number),
			zap.Int("page", opts.Page),
			zap.Int("perPage", opts.PerPage),
		)

		commits, _, err := r.client.PullRequests.ListCommits(ctx, organization, repository, number, &opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pull pull request commits")
		}

		if len(commits) == 0 || len(commits) < commitsPerPage {
			shouldRun = false
		}

		opts.Page++

		result = append(result, commits...)
	}

	return result, nil
}
