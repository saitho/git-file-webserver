package git

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/saitho/static-git-file-server/config"
)

const DownloadLocation = "./git_downloads"

type Client struct {
	Cfg *config.Config
}

func (c *Client) getCacheFilePath() string {
	return path.Join(c.getDownloadPath() + ".cache")
}

func (c *Client) getDownloadPath() string {
	return path.Join(DownloadLocation, c.Cfg.Git.Url)
}

func (c *Client) GetUpdatedTime() int64 {
	cacheFile, _ := ioutil.ReadFile(c.getCacheFilePath())
	cacheTime, _ := strconv.Atoi(string(cacheFile))
	return int64(cacheTime)
}

func (c *Client) IsUpToDate() bool {
	// File does not exist
	_, err := os.Stat(c.getDownloadPath())
	if os.IsNotExist(err) {
		return false
	}

	if c.Cfg.Git.Update.Mode != config.GitUpdateModeCache {
		return true
	}

	// Check if cache is up to date (within cacheTime interval)
	return time.Now().Unix() <= (c.GetUpdatedTime() + int64(c.Cfg.Git.Update.Cache.Time))
}
