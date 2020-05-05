package serialization

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/stefanoj3/gitstats/internal/domain/statistics"
)

func WritePullRequestStatistics(
	path string,
	stats statistics.PullRequestsStatistics,
) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "WritePullRequestStatistics: failed to create file `%s`", path)
	}

	defer file.Close() //nolint:errcheck

	csvWriter := csv.NewWriter(file)

	err = csvWriter.WriteAll(
		[][]string{
			{"TimeToMerge", "Merged", "Open", "Closed", "Total"},
			{
				stats.TimeToMerge.String(),
				strconv.Itoa(stats.Merged),
				strconv.Itoa(stats.Open),
				strconv.Itoa(stats.Closed),
				strconv.Itoa(stats.Total),
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "WritePullRequestStatistics: failed to write headers")
	}

	csvWriter.Flush()

	return nil
}
