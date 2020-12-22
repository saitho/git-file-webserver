package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"

	"github.com/saitho/static-git-file-server/config"
)

func (g *GitHandler) DownloadRepository() error {
	if err := os.RemoveAll(g.getDownloadPath()); err != nil {
		return fmt.Errorf("RemoveAll: %s", err.Error())
	}
	_, err := git.PlainClone(g.getDownloadPath(), false, &git.CloneOptions{
		URL:      g.Cfg.Git.Url,
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		return fmt.Errorf("PlainClone: %s", err.Error())
	}
	err = ioutil.WriteFile(g.getCacheFilePath(), []byte(strconv.Itoa(int(time.Now().Unix()))), 0644)
	if err != nil {
		return fmt.Errorf("WriteFile: %s", err.Error())
	}
	return nil
}

func (g *GitHandler) GetUpdatedTime() int64 {
	cacheFile, _ := ioutil.ReadFile(g.getCacheFilePath())
	cacheTime, _ := strconv.Atoi(string(cacheFile))
	return int64(cacheTime)
}

func (g *GitHandler) IsUpToDate() bool {
	// File does not exist
	_, err := os.Stat(g.getDownloadPath())
	if os.IsNotExist(err) {
		return false
	}

	if g.Cfg.Git.Update.Mode != config.GitUpdateModeCache {
		return true
	}

	// Check if cache is up to date (within cacheTime interval)
	return time.Now().Unix() <= (g.GetUpdatedTime() + int64(g.Cfg.Git.Update.Cache.Time))
}

func (g *GitHandler) GetBranches() []string {
	output, _ := g.runGitCommand("branch", "-l", "-r", "--no-color")
	var branches []string
	for _, v := range strings.Split(output, "\n") {
		v = strings.TrimSpace(v)
		v = strings.TrimPrefix(v, "origin/")
		branches = append(branches, v)
	}
	return branches
}

type GitTag struct {
	Tag  string
	Date time.Time
}

func (g *GitHandler) GetTags() []GitTag {
	sortPrefix := "-" // default: desc
	if strings.ToLower(g.Cfg.Display.Tags.Order) == "asc" {
		sortPrefix = ""
	}

	output, _ := g.runGitCommand("for-each-ref", "--sort="+sortPrefix+"creatordate", "--format=%(refname)---%(creatordate)", "refs/tags")
	var tags []GitTag
	for _, v := range strings.Split(output, "\n") {
		if v == "" {
			continue
		}
		split := strings.Split(v, "---")
		intDate, _ := time.Parse("Mon Jan 2 15:04:05 2006 -0700", split[1])
		tags = append(tags, GitTag{
			Tag:  strings.TrimPrefix(split[0], "refs/tags/"),
			Date: intDate,
		})
	}
	return tags
}
