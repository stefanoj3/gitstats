package git

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

const commitsPerPage = 100

func NewGithubCommitsRepository(client *github.Client) *GithubCommitsRepository {
	return &GithubCommitsRepository{client: client}
}

type GithubCommitsRepository struct {
	client *github.Client
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
	for shouldRun {
		fmt.Println("listing commits for", organization, repository, number, &opts)

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
