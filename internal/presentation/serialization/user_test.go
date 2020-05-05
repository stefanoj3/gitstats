package serialization_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stefanoj3/gitstats/internal/domain/statistics"
	"github.com/stefanoj3/gitstats/internal/presentation/serialization"
	"github.com/stretchr/testify/assert"
)

func TestWriteUsersStatistics(t *testing.T) {
	const (
		filePath = "testdata/out/user_stats.csv"

		user1 = "user_1"
		user2 = "user_2"
	)

	basetime := time.Date(2007, 7, 7, 0, 0, 0, 0, time.UTC)

	stats := statistics.NewUserStatistics()
	stats.At(user1, basetime).Commits += 2
	stats.At(user1, basetime).Comments += 10
	stats.At(user1, basetime).PullRequestsCreated = 1

	stats.At(user1, basetime.AddDate(0, 0, 10)).Commits += 2
	stats.At(user1, basetime.AddDate(0, 0, 10)).Comments += 10
	stats.At(user1, basetime.AddDate(0, 0, 10)).PullRequestsCreated = 1

	stats.At(user2, basetime).Commits += 5
	stats.At(user2, basetime).Comments += 15
	stats.At(user2, basetime).PullRequestsCreated = 2

	defer os.Remove(filePath) //nolint:errcheck

	err := serialization.WriteUsersStatistics(
		filePath,
		stats,
		basetime.Add(-time.Hour*24),
		basetime.Add(time.Hour*24),
		[]string{user1},
	)

	assert.NoError(t, err)

	b, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	content := string(b)

	expectedContent := `Date,User,PullRequestsCreated,Comments,Commits
2007-07-06,user_1,0,0,0
2007-07-07,user_1,1,10,2
2007-07-08,user_1,0,0,0
`

	assert.Equal(t, expectedContent, content)
}

func TestWriteUsersStatisticsShouldFailWithInvalidPath(t *testing.T) {
	const (
		filePath = "/root/something/prs.csv"

		user1 = "user_1"
	)

	basetime := time.Date(2007, 7, 7, 0, 0, 0, 0, time.UTC)

	stats := statistics.NewUserStatistics()

	err := serialization.WriteUsersStatistics(
		filePath,
		stats,
		basetime.Add(-time.Hour*24),
		basetime.Add(time.Hour*24),
		[]string{user1},
	)
	assert.Error(t, err)
}
