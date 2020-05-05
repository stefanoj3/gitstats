package cmd

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/domain/statistics"
	"github.com/stefanoj3/gitstats/internal/infrastructure/git"
	"github.com/stefanoj3/gitstats/internal/infrastructure/oauth"
	"github.com/stefanoj3/gitstats/internal/presentation/serialization"
	"github.com/stefanoj3/gitstats/internal/usecase"
	"go.uber.org/zap"
)

func NewCollectCommand(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect statistics",
		RunE:  buildCollectCommand(logger),
	}

	cmd.Flags().StringP(
		flagCollectConfigFile,
		flagCollectConfigFileShort,
		"",
		"configuration file for your team",
	)
	Must(cmd.MarkFlagFilename(flagCollectConfigFile))
	Must(cmd.MarkFlagRequired(flagCollectConfigFile))

	cmd.Flags().StringP(
		flagCollectFromDate,
		flagCollectFromDateShort,
		"",
		"from when we need to start collecting stats (Y-m-d format)",
	)
	Must(cmd.MarkFlagRequired(flagCollectFromDate))

	cmd.Flags().DurationP(
		flagCollectDelta,
		flagCollectDeltaShort,
		time.Hour*24*5,
		"delta time used to search for comments/commits between the from/to in PRs in range bigger than the specific from/to",
	)

	cmd.Flags().StringP(
		flagCollectToDate,
		flagCollectToDateShort,
		"",
		"to when we need to stop collecting stats (Y-m-d format)",
	)
	Must(cmd.MarkFlagRequired(flagCollectToDate))

	cmd.Flags().StringP(
		flagOutputFilePrefix,
		flagOutputFilePrefixShort,
		"out",
		"prefix for the output file",
	)

	return cmd
}

func buildCollectCommand(logger *zap.Logger) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger.Info("Initializing")

		ctx := context.Background()

		token, err := getGithubToken()
		if err != nil {
			return errors.Wrap(err, "CollectCommand: failed to get token")
		}

		config, err := getCollectConfig(cmd)
		if err != nil {
			return errors.Wrap(err, "CollectCommand: failed to get config")
		}

		client := github.NewClient(
			oauth.NewClient(ctx, token),
		)

		githubAggregatedRepository := git.NewGithubAggregatedRepository(
			git.NewGithubPullRequestsRepository(client, logger),
			git.NewGithubCommitsRepository(client, logger),
			git.NewGithubCommentsRepository(client, logger),
			logger,
		)

		getStatisticsUseCase := usecase.NewGetStatistics(
			githubAggregatedRepository,
			logger,
		)

		logger.Info(
			"Starting to collect statistics",
			zap.String("organization", config.Organization),
			zap.Strings("repositories", config.Repositories),
			zap.Strings("users", config.Users),
			zap.Time("from", config.From),
			zap.Time("to", config.To),
			zap.Duration("delta", config.Delta),
			zap.Int("tokenLen", len(token)),
		)

		stats, err := getStatisticsUseCase.GetStatistics(
			ctx,
			config.From,
			config.To,
			config.Delta,
			config.Organization,
			config.Repositories,
			config.Users,
		)
		if err != nil {
			return errors.Wrap(err, "CollectCommand: failed to get statistics")
		}

		logger.Info(
			"Done",
			zap.Duration("timeToMerge", stats.PullRequestsStatistics.TimeToMerge),
			zap.Int("merged", stats.PullRequestsStatistics.Merged),
			zap.Int("open", stats.PullRequestsStatistics.Open),
			zap.Int("closed", stats.PullRequestsStatistics.Closed),
			zap.Int("total", stats.PullRequestsStatistics.Total),
		)

		return writeOutput(config, stats, logger)
	}
}

func writeOutput(config CollectConfig, stats *statistics.Statistics, logger *zap.Logger) error {
	pullRequestsOutFile := config.OutputFilePrefix + "_pull_requests.csv"

	logger.Info("Writing output file for pull requests", zap.String("path", pullRequestsOutFile))

	err := serialization.WritePullRequestStatistics(pullRequestsOutFile, stats.PullRequestsStatistics)
	if err != nil {
		return errors.Wrap(err, "CollectCommand: failed to write pull requests statistics")
	}

	usersOutFile := config.OutputFilePrefix + "_user.csv"

	logger.Info("Writing output file for users", zap.String("path", usersOutFile))

	err = serialization.WriteUsersStatistics(
		usersOutFile,
		stats.UsersStatistics,
		config.From,
		config.To,
		config.Users,
	)
	if err != nil {
		return errors.Wrap(err, "CollectCommand: failed to write users statistics")
	}

	return nil
}

func getGithubToken() (string, error) {
	// nolint:gosec
	const githubTokenEnvVariable = "GITHUB_TOKEN"

	token := os.Getenv(githubTokenEnvVariable)
	if len(token) == 0 {
		return "", errors.Errorf(
			"CollectCommand: please provide a github token in the environment (%s)",
			githubTokenEnvVariable,
		)
	}

	return token, nil
}
