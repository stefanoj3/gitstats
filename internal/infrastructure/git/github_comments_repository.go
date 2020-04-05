package git

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func NewGithubCommentsRepository(client *github.Client) *GithubCommentsRepository {
	return &GithubCommentsRepository{client: client}
}

type GithubCommentsRepository struct {
	client *github.Client
}

func (r *GithubCommentsRepository) FetchCommentsFor(
	ctx context.Context,
	organization string,
	repository string,
	number int,
	usersHandles []string,
) ([]*github.PullRequestComment, error) {
	opts := github.PullRequestListCommentsOptions{
		Sort:        "created",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
	}

	result := make([]*github.PullRequestComment, 0, 100)

	for {
		comments, _, err := r.client.PullRequests.ListComments(ctx, organization, repository, number, &opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch pull pull request comments")
		}

		if len(comments) == 0 {
			break
		}

		for _, comment := range comments {
			if !userIn(*comment.User.Login, usersHandles) {
				continue
			}

			result = append(result, comment)
		}
	}

	return result, nil
}
