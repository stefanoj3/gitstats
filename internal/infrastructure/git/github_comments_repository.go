package git

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const commentsPerPage = 50

func NewGithubCommentsRepository(client *github.Client) *GithubCommentsRepository {
	return &GithubCommentsRepository{client: client}
}

type GithubCommentsRepository struct {
	client *github.Client
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
	for shouldRun {
		fmt.Println("listing comments for", organization, repository, number, &opts)

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
