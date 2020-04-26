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
}
