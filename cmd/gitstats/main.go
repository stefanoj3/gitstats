package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/presentation/cli/cmd"
)

func main() {
	if err := buildCLICommand().Execute(); err != nil {
		log.Fatalln("failed to run", err)
	}
}

func buildCLICommand() *cobra.Command {
	root := cmd.NewRootCommand()

	root.AddCommand(cmd.NewCollectCommand())
	root.AddCommand(cmd.NewVersionCommand(os.Stdout))

	return root
}
