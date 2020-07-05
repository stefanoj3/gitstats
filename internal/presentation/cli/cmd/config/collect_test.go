package config_test

import (
	"testing"

	"github.com/stefanoj3/gitstats/internal/presentation/cli/cmd/config"
	"github.com/stretchr/testify/assert"
)

func TestCollectConfig_Merge(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		configs        []config.CollectConfig
		expectedResult config.CollectConfig
	}{
		{
			scenario: "one empty one filled",
			configs: []config.CollectConfig{
				{},
				{
					Organization: "myorg",
					Repositories: []string{"myrepo1", "myrepo2"},
					Users:        []string{"user1", "user2"},
				},
			},
			expectedResult: config.CollectConfig{
				Organization: "myorg",
				Repositories: []string{"myrepo1", "myrepo2"},
				Users:        []string{"user1", "user2"},
			},
		},
		{
			scenario: "organization of the first is overridden by the second",
			configs: []config.CollectConfig{
				{
					Organization: "myorg1",
				},
				{
					Organization: "myorg2",
					Repositories: []string{"myrepo1", "myrepo2"},
					Users:        []string{"user1", "user2"},
				},
			},
			expectedResult: config.CollectConfig{
				Organization: "myorg2",
				Repositories: []string{"myrepo1", "myrepo2"},
				Users:        []string{"user1", "user2"},
			},
		},
		{
			scenario: "repositories and users are merged",
			configs: []config.CollectConfig{
				{
					Repositories: []string{"myrepo1", "myrepo2"},
					Users:        []string{"user1", "user2"},
				},
				{
					Organization: "myorg1",
					Repositories: []string{"myrepo3", "myrepo4"},
					Users:        []string{"user3", "user4"},
				},
			},
			expectedResult: config.CollectConfig{
				Organization: "myorg1",
				Repositories: []string{"myrepo1", "myrepo2", "myrepo3", "myrepo4"},
				Users:        []string{"user1", "user2", "user3", "user4"},
			},
		},
		{
			scenario: "duplicated users and repos are removed",
			configs: []config.CollectConfig{
				{
					Repositories: []string{"myrepo1", "myrepo2", "myrepo3"},
					Users:        []string{"user1", "user2", "user3"},
				},
				{
					Organization: "myorg1",
					Repositories: []string{"myrepo3", "myrepo4"},
					Users:        []string{"user3", "user4"},
				},
			},
			expectedResult: config.CollectConfig{
				Organization: "myorg1",
				Repositories: []string{"myrepo1", "myrepo2", "myrepo3", "myrepo4"},
				Users:        []string{"user1", "user2", "user3", "user4"},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint

		t.Run(tc.scenario, func(t *testing.T) {
			var result config.CollectConfig

			for _, c := range tc.configs {
				result = result.Merge(c)
			}

			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
