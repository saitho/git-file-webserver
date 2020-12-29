package git

import (
	"fmt"
	"strings"

	"github.com/saitho/static-git-file-server/config"
)

type Client struct {
	Cfg         *config.Config
	CurrentRepo *config.RepoConfig
}

func (c *Client) GetRepositoryBySlug(slug string) *config.RepoConfig {
	for _, repository := range c.Cfg.Git.Repositories {
		if strings.ToLower(repository.Slug) == strings.ToLower(slug) {
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
