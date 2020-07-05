package cmd

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/stefanoj3/gitstats/internal/presentation/cli/cmd/config"
)

const timeFormatLayout = "2006-01-02"
const failedToParseTimeErrorString = "getCollectConfig: failed to parse `%s`, expected format is `Y-m-d`"

func getCollectConfig(cmd *cobra.Command) (config.CollectConfig, error) {
	var c config.CollectConfig

	configFilePaths, err := cmd.Flags().GetStringSlice(flagCollectConfigFile)
	if err != nil {
		return c, errors.Wrapf(err, "failed to read %s", flagCollectConfigFile)
	}

	for _, configFilePath := range configFilePaths {
		var temporaryConfig config.CollectConfig

		_, err = toml.DecodeFile(configFilePath, &temporaryConfig)
		if err != nil {
			return c, errors.Wrapf(err, "getCollectConfig: failed to get config from %s", configFilePath)
		}

		c = c.Merge(temporaryConfig)
	}

	rawFromDate := cmd.Flag(flagCollectFromDate).Value.String()

	c.From, err = time.Parse(timeFormatLayout, rawFromDate)
	if err != nil {
		return c, errors.Wrapf(err, failedToParseTimeErrorString, flagCollectFromDate)
	}

	rawToDate := cmd.Flag(flagCollectToDate).Value.String()

	c.To, err = time.Parse(timeFormatLayout, rawToDate)
	if err != nil {
		return c, errors.Wrapf(err, failedToParseTimeErrorString, flagCollectToDate)
	}

	c.Delta, err = cmd.Flags().GetDuration(flagCollectDelta)
	if err != nil {
		return c, errors.Wrapf(err, "failed to parse %s", flagCollectDelta)
	}

	c.OutputFilePrefix = cmd.Flag(flagOutputFilePrefix).Value.String()

	if len(c.Organization) == 0 {
		return c, errors.New("invalid config: no `Organization` specified")
	}

	return c, nil
}
