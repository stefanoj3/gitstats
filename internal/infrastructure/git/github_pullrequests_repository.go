package git

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func NewGithubPullRequestsRepository(client *github.Client) *GithubPullRequestsRepository {
	return &GithubPullRequestsRepository{client: client}
}

type GithubPullRequestsRepository struct {
	client *github.Client
}

func (r *GithubPullRequestsRepository) FetchPullRequestsFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	organization string,
	repositories []string,
	usersHandles []string,
) ([]*github.PullRequest, error) {
	result := make([]*github.PullRequest, 0)

	for _, repository := range repositories {
		prs, err := r.fetchAllFor(ctx, from, to, organization, repository, usersHandles)
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
	usersHandles []string,
) ([]*github.PullRequest, error) {
	listOptions := github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		State:     "all",
		Sort:      "created",
		Direction: "desc",
	}

	fmt.Println("performing call for", repository, listOptions)

	result := make([]*github.PullRequest, 0, 100)

	shouldRun := true
	for shouldRun {
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

			if !userIn(*pr.User.Login, usersHandles) {
				continue
			}

			result = append(result, pr)
		}

		listOptions.Page++
	}

	return result, nil
}
