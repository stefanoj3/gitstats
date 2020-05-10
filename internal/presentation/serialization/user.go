package serialization

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/stefanoj3/gitstats/internal/domain/statistics"
)

const defaultTimeFormat = "2006-01-02"

func WriteUsersStatistics(
	path string,
	stats statistics.UsersStatistics,
	from, to time.Time,
	users []string,
) error {
	rangedTime := calculateRange(from, to)

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file `%s`", path)
	}

	defer file.Close() //nolint:errcheck

	csvWriter := csv.NewWriter(file)

	err = csvWriter.Write([]string{"Date", "User", "PullRequestsCreated", "Comments", "Commits"})
	if err != nil {
		return errors.Wrap(err, "WriteUsersStatistics: failed to write headers")
	}

	for _, user := range users {
		for _, t := range rangedTime {
			s := stats.At(user, t)

			err = csvWriter.Write([]string{
				t.Format(defaultTimeFormat),
				user,
				strconv.Itoa(s.PullRequestsCreated),
				strconv.Itoa(s.Comments),
				strconv.Itoa(s.Commits),
			})
			if err != nil {
				return errors.Wrapf(err, "WriteUsersStatistics: failed to write row for %s/%s", user, t.Format(time.RFC3339))
			}
		}
	}

	csvWriter.Flush()

	return nil
}

func calculateRange(from time.Time, to time.Time) []time.Time {
	rangedTime := make([]time.Time, 0, 30)
	start := from
	rangedTime = append(rangedTime, start)

	for {
		if start.Before(to) {
			start = start.Add(time.Hour * 24)
		} else {
			break
		}

		rangedTime = append(rangedTime, start)
	}

	return rangedTime
}
