package git

import (
	"io/ioutil"
	"os"
	"strconv"
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
	if time.Now().Unix() > int64(cacheTime+g.Cfg.Git.CacheTime) {
		return false
	}
	return true
}
