package git

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func (g *GitHandler) isUpToDate() bool {
	// File does not exist
	_, err := os.Stat(g.getDownloadPath())
	if os.IsNotExist(err) {
		return false
	}

	// Check download date
	cacheFile, _ := ioutil.ReadFile(g.getCacheFilePath())
	cacheTime, _ := strconv.Atoi(string(cacheFile))
	return time.Now().Unix() <= int64(cacheTime+g.Cfg.Git.CacheTime)
}

func (g *GitHandler) GetBranches() []string {
	output, _ := g.runGitCommand("branch", "-l", "--no-color")
	var branches []string
	for _, v := range strings.Split(output, "\n") {
		v = strings.TrimPrefix(v, "*")
		v = strings.TrimSpace(v)
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
