package cmd

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const timeFormatLayout = "2006-01-02"
const failedToParseTimeErrorString = "getCollectConfig: failed to parse `%s`, expected format is `Y-m-d`"

type CollectConfig struct {
	Organization string
	Repositories []string
	// Users should contain the handles of the users
	Users []string
	From  time.Time
	To    time.Time
}

func getCollectConfig(cmd *cobra.Command) (CollectConfig, error) {
	var (
		config CollectConfig
		err    error
	)

	configFilePath := cmd.Flag(flagCollectConfigFile).Value.String()

	_, err = toml.DecodeFile(configFilePath, &config)
	if err != nil {
		return config, errors.Wrapf(err, "getCollectConfig: failed to get config from %s", configFilePath)
	}

	rawFromDate := cmd.Flag(flagCollectFromDate).Value.String()

	config.From, err = time.Parse(timeFormatLayout, rawFromDate)
	if err != nil {
		return config, errors.Wrapf(err, failedToParseTimeErrorString, flagCollectFromDate)
	}

	rawToDate := cmd.Flag(flagCollectToDate).Value.String()

	config.To, err = time.Parse(timeFormatLayout, rawToDate)
	if err != nil {
		return config, errors.Wrapf(err, failedToParseTimeErrorString, flagCollectToDate)
	}

	return config, nil
}
