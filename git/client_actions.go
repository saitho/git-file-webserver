package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"

	"github.com/saitho/static-git-file-server/config"
)

func (c *Client) runGitCommand(repo *config.RepoConfig, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.GetDownloadPath()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *Client) DownloadRepository(repo *config.RepoConfig) error {
	if err := os.RemoveAll(repo.GetDownloadPath()); err != nil {
		return fmt.Errorf("RemoveAll: %s", err.Error())
	}
	_, err := git.PlainClone(repo.GetDownloadPath(), false, &git.CloneOptions{
		URL:      repo.Url,
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return fmt.Errorf("PlainClone: %s", err.Error())
	}
	err = ioutil.WriteFile(repo.GetCacheFilePath(), []byte(strconv.Itoa(int(time.Now().Unix()))), 0644)
	if err != nil {
		return fmt.Errorf("WriteFile: %s", err.Error())
	}
	return nil
}
