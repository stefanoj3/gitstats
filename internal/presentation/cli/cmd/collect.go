package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/infrastructure/git"
	"github.com/stefanoj3/gitstats/internal/infrastructure/oauth"
	"github.com/stefanoj3/gitstats/internal/usecase"
)

func NewCollectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect statistics",
		RunE:  collectCommand,
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

	cmd.Flags().StringP(
		flagCollectToDate,
		flagCollectToDateShort,
		"",
		"to when we need to stop collecting stats (Y-m-d format)",
	)
	Must(cmd.MarkFlagRequired(flagCollectToDate))

	return cmd
}

func collectCommand(cmd *cobra.Command, args []string) error {
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

	getStatisticsUseCase := usecase.NewGetStatistics(
		git.NewGithubPullRequestsRepository(client),
		git.NewGithubCommentsRepository(client),
	)

	stats, err := getStatisticsUseCase.GetStatistics(
		ctx,
		config.From,
		config.To,
		config.Organization,
		config.Repositories,
		config.Users,
	)
	if err != nil {
		return errors.Wrap(err, "CollectCommand: failed to get statistics")
	}

	fmt.Println(stats)

	return nil
}

func getGithubToken() (string, error) {
	// nolint:gosec
	const githubTokenEnvVariable = "GITHUB_TOKEN"

	token := os.Getenv(githubTokenEnvVariable)
	if len(token) == 0 {
		return "", errors.New(
			fmt.Sprintf(
				"CollectCommand: please provide a github token in the environment (%s)",
				githubTokenEnvVariable,
			),
		)
	}

	return token, nil
}
