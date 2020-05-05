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

func TestWritePullRequestStatistics(t *testing.T) {
	const filePath = "testdata/out/prs.csv"

	stats := statistics.PullRequestsStatistics{
		TimeToMerge: time.Second,
		Merged:      1,
		Closed:      2,
		Open:        3,
		Total:       6,
	}

	defer os.Remove(filePath) //nolint:errcheck

	err := serialization.WritePullRequestStatistics(filePath, stats)
	assert.NoError(t, err)

	b, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)

	content := string(b)

	expectedContent := `TimeToMerge,Merged,Open,Closed,Total
1s,1,3,2,6
`
	assert.Equal(t, expectedContent, content)
}

func TestWritePullRequestStatisticsShouldFailWithInvalidPath(t *testing.T) {
	const filePath = "/root/something/prs.csv"

	stats := statistics.PullRequestsStatistics{
		TimeToMerge: time.Second,
		Merged:      1,
		Closed:      2,
		Open:        3,
		Total:       6,
	}

	err := serialization.WritePullRequestStatistics(filePath, stats)
	assert.Error(t, err)
}
