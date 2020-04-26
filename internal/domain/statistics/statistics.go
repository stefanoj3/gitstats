package statistics

import "time"

type DailyStatistics struct {
	Commits             int
	Comments            int
	PullRequestsCreated int
}

func NewUserStatistics() UsersStatistics {
	return UsersStatistics{dailyStatsByLoginAndDate: map[string]map[time.Time]*DailyStatistics{}}
}

type UsersStatistics struct {
	dailyStatsByLoginAndDate map[string]map[time.Time]*DailyStatistics
}

func (s *UsersStatistics) At(login string, t time.Time) *DailyStatistics {
	dailyStatsByDate, ok := s.dailyStatsByLoginAndDate[login]
	if !ok {
		s.dailyStatsByLoginAndDate[login] = make(map[time.Time]*DailyStatistics)
		dailyStatsByDate = s.dailyStatsByLoginAndDate[login]
	}

	y, m, d := t.Date()
	t = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	_, ok = dailyStatsByDate[t]
	if !ok {
		dailyStatsByDate[t] = &DailyStatistics{}
	}

	return dailyStatsByDate[t]
}

type PullRequestsStatistics struct {
	// TimeToMerge represents the average time it takes for a PR to be merged
	TimeToMerge time.Duration
	// Merged represents the total or PRs that got merged
	Merged int
	// Closed represents the total or PRs that got closed but never merged
	Closed int
	// Open represents the total or PRs that are still open
	Open int
	// Total is the total amount of PRs
	Total int
}

type Statistics struct {
	PullRequestsStatistics PullRequestsStatistics
	UsersStatistics        UsersStatistics
}
