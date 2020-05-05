package git

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const commentsPerPage = 50

func NewGithubCommentsRepository(client *github.Client, logger *zap.Logger) *GithubCommentsRepository {
	return &GithubCommentsRepository{client: client, logger: logger}
}

type GithubCommentsRepository struct {
	client *github.Client
	logger *zap.Logger
}

func (r *GithubCommentsRepository) FindCommentsFor(
	ctx context.Context,
	organization string,
	repository string,
	number int,
) ([]*github.PullRequestComment, error) {
	opts := github.PullRequestListCommentsOptions{
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: commentsPerPage, Page: 1},
	}

	result := make([]*github.PullRequestComment, 0, commentsPerPage)

	shouldRun := true

	//nolint:dupl
	for shouldRun {
		r.logger.Debug(
			"listing comments",
			zap.String("organization", organization),
			zap.String("repository", organization),
			zap.Int("number", number),
			zap.Int("page", opts.Page),
			zap.Int("perPage", opts.PerPage),
		)

		comments, _, err := r.client.PullRequests.ListComments(ctx, organization, repository, number, &opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pull pull request comments")
		}

		if len(comments) == 0 || len(comments) < commentsPerPage {
			shouldRun = false
		}

		opts.Page++

		result = append(result, comments...)
	}

	return result, nil
}
