package cmd

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitstats",
		Short: "TODO write short description",
		Long:  "TODO write long description",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
	}

	return cmd
}
