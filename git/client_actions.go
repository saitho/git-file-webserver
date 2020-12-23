package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-git/go-git/v5"
)

func (c *Client) runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = c.getDownloadPath()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *Client) DownloadRepository() error {
	if err := os.RemoveAll(c.getDownloadPath()); err != nil {
		return fmt.Errorf("RemoveAll: %s", err.Error())
	}
	_, err := git.PlainClone(c.getDownloadPath(), false, &git.CloneOptions{
		URL:      c.Cfg.Git.Url,
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return fmt.Errorf("PlainClone: %s", err.Error())
	}
	err = ioutil.WriteFile(c.getCacheFilePath(), []byte(strconv.Itoa(int(time.Now().Unix()))), 0644)
	if err != nil {
		return fmt.Errorf("WriteFile: %s", err.Error())
	}
	return nil
}
