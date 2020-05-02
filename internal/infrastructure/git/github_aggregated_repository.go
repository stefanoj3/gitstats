package git

import (
	"context"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

func NewGithubAggregatedRepository(
	pullRequestsRepository *GithubPullRequestsRepository,
	commitsRepository *GithubCommitsRepository,
	commentRepository *GithubCommentsRepository,
) *GithubAggregatedRepository {
	return &GithubAggregatedRepository{pullRequestsRepository: pullRequestsRepository, commitsRepository: commitsRepository, commentRepository: commentRepository}
}

type GithubAggregatedRepository struct {
	pullRequestsRepository *GithubPullRequestsRepository
	commitsRepository      *GithubCommitsRepository
	commentRepository      *GithubCommentsRepository
}

func (r *GithubAggregatedRepository) FetchAllFor(
	ctx context.Context,
	from time.Time,
	to time.Time,
	delta time.Duration,
	organization string,
	repositories []string,
) ([]*github.PullRequest, []*github.RepositoryCommit, []*github.PullRequestComment, error) {
	pullRequests, err := r.pullRequestsRepository.FindPullRequestsFor(
		ctx,
		from.Add(-delta),
		to.Add(delta),
		organization,
		repositories,
	)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "FetchAllFor: failed to get pull requests")
	}

	commits, err := r.getCommitsFor(ctx, pullRequests)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "FetchAllFor: failed to get commits")
	}

	comments, err := r.getCommentsFor(ctx, pullRequests)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "FetchAllFor: failed to get comments")
	}

	return filterPullRequestsByDateRange(from, to, pullRequests),
		filterCommitsByDateRange(from, to, commits),
		filterCommentsByDateRange(from, to, comments),
		nil
}

func filterCommentsByDateRange(from, to time.Time, comments []*github.PullRequestComment) []*github.PullRequestComment {
	result := make([]*github.PullRequestComment, 0, len(comments))

	for _, comment := range comments {
		if comment.CreatedAt.Before(from) {
			continue
		}

		if comment.CreatedAt.After(to) {
			continue
		}

		result = append(result, comment)
	}

	return result
}

func filterCommitsByDateRange(from, to time.Time, commits []*github.RepositoryCommit) []*github.RepositoryCommit {
	result := make([]*github.RepositoryCommit, 0, len(commits))

	for _, commit := range commits {
		if commit.Commit.Committer.Date.Before(from) {
			continue
		}

		if commit.Commit.Committer.Date.After(to) {
			continue
		}

		result = append(result, commit)
	}

	return result
}

func filterPullRequestsByDateRange(from, to time.Time, pullRequests []*github.PullRequest) []*github.PullRequest {
	result := make([]*github.PullRequest, 0, len(pullRequests))

	for _, pr := range pullRequests {
		if pr.CreatedAt.Before(from) {
			continue
		}

		if pr.CreatedAt.After(to) {
			continue
		}

		result = append(result, pr)
	}

	return result
}

// nolint:dupl
func (r *GithubAggregatedRepository) getCommentsFor(
	ctx context.Context,
	pullRequests []*github.PullRequest,
) ([]*github.PullRequestComment, error) {
	result := make([]*github.PullRequestComment, 0, len(pullRequests))

	for _, pr := range pullRequests {
		comments, err := r.commentRepository.FindCommentsFor(
			ctx,
			*pr.Base.Repo.Owner.Login,
			*pr.Base.Repo.Name,
			*pr.Number,
		)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"failed to fetch comments for %s/%s/%d",
				*pr.Base.Repo.Owner.Login,
				*pr.Base.Repo.Name,
				*pr.Number,
			)
		}

		result = append(result, comments...)
	}

	return result, nil
}

// nolint:dupl
func (r *GithubAggregatedRepository) getCommitsFor(
	ctx context.Context,
	pullRequests []*github.PullRequest,
) ([]*github.RepositoryCommit, error) {
	result := make([]*github.RepositoryCommit, 0, len(pullRequests))

	for _, pr := range pullRequests {
		commits, err := r.commitsRepository.FindCommitsFor(
			ctx,
			*pr.Base.Repo.Owner.Login,
			*pr.Base.Repo.Name,
			*pr.Number,
		)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"failed to fetch commits for %s/%s/%d",
				*pr.Base.Repo.Owner.Login,
				*pr.Base.Repo.Name,
				*pr.Number,
			)
		}

		result = append(result, commits...)
	}

	return result, nil
}
