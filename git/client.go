package git

import (
	"fmt"
	"strings"

	"github.com/saitho/static-git-file-server/config"
)

type ClientInterface interface {
	GetRepositoryBySlug(slug string) *config.RepoConfig
	SelectRepository(slug string) error
	GetTags(repo *config.RepoConfig) []GitTag
	GetCurrentRepo() *config.RepoConfig
}

type Client struct {
	Cfg         *config.Config
	CurrentRepo *config.RepoConfig
}

func (c *Client) GetCurrentRepo() *config.RepoConfig {
	return c.CurrentRepo
}

func (c *Client) GetRepositoryBySlug(slug string) *config.RepoConfig {
	for _, repository := range c.Cfg.Git.Repositories {
		if strings.EqualFold(repository.Slug, slug) {
			return repository
		}
	}
	return nil
}

func (c *Client) SelectRepository(slug string) error {
	if len(slug) == 0 {
		c.CurrentRepo = nil
		return nil
	}
	c.CurrentRepo = c.GetRepositoryBySlug(slug)
	if c.CurrentRepo == nil {
		return fmt.Errorf("unable to select repository with slug \"%s\"", slug)
	}
	return nil
}
