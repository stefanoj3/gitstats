package statistics_test

import (
	"testing"
	"time"

	"github.com/stefanoj3/gitstats/internal/domain/statistics"
	"github.com/stretchr/testify/assert"
)

var baseTime = time.Date(2007, time.July, 7, 0, 0, 0, 0, time.UTC)
const (
	login1 = "test-user1"
	login2 = "test-user2"
)

func TestUserStatistics(t *testing.T) {
	stats := statistics.NewUserStatistics()

	stats.At(login1, baseTime).Commits++
	assert.Equal(t, 1, stats.At(login1, baseTime).Commits)
	assert.Equal(t, 0, stats.At(login2, baseTime).Commits)

	stats.At(login1, baseTime).Comments++
	stats.At(login1, baseTime).Comments++
	assert.Equal(t, 2, stats.At(login1, baseTime).Comments)
	assert.Equal(t, 0, stats.At(login2, baseTime).Comments)

	stats.At(login1, baseTime).PullRequestsCreated++
	stats.At(login1, baseTime).PullRequestsCreated++
	stats.At(login1, baseTime).PullRequestsCreated++
	stats.At(login1, baseTime).PullRequestsCreated++
	assert.Equal(t, 4, stats.At(login1, baseTime).PullRequestsCreated)
	assert.Equal(t, 0, stats.At(login2, baseTime).PullRequestsCreated)

	afterBaseTime := baseTime.AddDate(0, 0, 1)

	stats.At(login1, afterBaseTime).PullRequestsCreated++
	assert.Equal(t, 1, stats.At(login1, afterBaseTime).PullRequestsCreated)
	assert.Equal(t, 0, stats.At(login2, afterBaseTime).PullRequestsCreated)

	times := make([]time.Time, 0, 2)
	result := make([]*statistics.DailyStatistics, 0, 2)
	stats.Range(baseTime, afterBaseTime, func(login string, t time.Time, dailyStatistics *statistics.DailyStatistics) {
		if login != login1 {
			return
		}

		times = append(times, t)
		result = append(result, dailyStatistics)
	})

	assert.Len(t, times, 2)
	assert.Equal(t, baseTime, times[0])
	assert.Equal(t, afterBaseTime, times[1])

	assert.Equal(t, 4, result[0].PullRequestsCreated)
	assert.Equal(t, 1, result[1].PullRequestsCreated)
}
