package config

import (
	"sort"
	"time"
)

type CollectConfig struct {
	Organization string
	Repositories []string
	// Users should contain the handles of the users
	Users            []string
	From             time.Time
	To               time.Time
	Delta            time.Duration
	OutputFilePrefix string
}

func (c CollectConfig) Merge(other CollectConfig) CollectConfig {
	if len(other.Organization) > 0 {
		c.Organization = other.Organization
	}

	repositories := c.mergeRepositories(other)
	c.Repositories = make([]string, 0, len(repositories))

	for r := range repositories {
		c.Repositories = append(c.Repositories, r)
	}

	users := c.mergeUsers(other)
	c.Users = make([]string, 0, len(users))

	for r := range users {
		c.Users = append(c.Users, r)
	}

	sort.Strings(c.Repositories)
	sort.Strings(c.Users)

	return c
}

func (c CollectConfig) mergeUsers(other CollectConfig) map[string]interface{} {
	users := make(map[string]interface{}, len(c.Users)+len(other.Users))

	for _, r := range c.Users {
		users[r] = nil
	}

	for _, r := range other.Users {
		users[r] = nil
	}

	return users
}

func (c CollectConfig) mergeRepositories(other CollectConfig) map[string]interface{} {
	repositories := make(map[string]interface{}, len(c.Repositories)+len(other.Repositories))

	for _, r := range c.Repositories {
		repositories[r] = nil
	}

	for _, r := range other.Repositories {
		repositories[r] = nil
	}

	return repositories
}
