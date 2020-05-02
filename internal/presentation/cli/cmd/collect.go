package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/domain/statistics"
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

	githubAggregatedRepository := git.NewGithubAggregatedRepository(
		git.NewGithubPullRequestsRepository(client),
		git.NewGithubCommitsRepository(client),
		git.NewGithubCommentsRepository(client),
	)

	getStatisticsUseCase := usecase.NewGetStatistics(
		githubAggregatedRepository,
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

	printResults(stats, config)

	return nil
}

func printResults(stats *statistics.Statistics, config CollectConfig) {
	rangedTime := make([]time.Time, 0)
	start := config.From
	rangedTime = append(rangedTime, start)

	for {
		if start.Before(config.To) {
			start = start.Add(time.Hour * 24)
		} else {
			break
		}

		rangedTime = append(rangedTime, start)
	}

	for _, user := range config.Users {
		for _, t := range rangedTime {
			s := stats.UsersStatistics.At(user, t)
			fmt.Println(
				fmt.Sprintf("%s login: %s - PRs: %d, Comments: %d, Commits: %d", t.Format(time.RFC3339), user, s.PullRequestsCreated, s.Comments, s.Commits))
		}
	}

	fmt.Println("TimeToMerge:", stats.PullRequestsStatistics.TimeToMerge)
	fmt.Println("Merged:", stats.PullRequestsStatistics.Merged)
	fmt.Println("Closed:", stats.PullRequestsStatistics.Closed)
	fmt.Println("Open:", stats.PullRequestsStatistics.Open)
	fmt.Println("Total:", stats.PullRequestsStatistics.Total)
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
