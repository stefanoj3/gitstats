package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewRootCommand(zapAtomicLevel *zap.AtomicLevel) *cobra.Command {
	var verbose bool

	cmd := &cobra.Command{
		Use:   "gitstats",
		Short: "gitstats is a software to fetch statistics about your team activity from github",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				zapAtomicLevel.SetLevel(zap.DebugLevel)
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(
		&verbose,
		flagRootVerbose,
		flagRootVerboseShort,
		false,
		"verbose mode",
	)

	return cmd
}
