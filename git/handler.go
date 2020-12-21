package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/saitho/static-git-file-server/config"
)

const DownloadLocation = "./git_downloads"

type GitHandler struct {
	Cfg config.Config
}

func (g *GitHandler) getCacheFilePath() string {
	return path.Join(g.getDownloadPath() + ".cache")
}
func (g *GitHandler) getDownloadPath() string {
	return path.Join(DownloadLocation, g.Cfg.Git.Url)
}

func (g *GitHandler) ServePath(refType string, refName string, filePath string) (string, error) {
	if !g.isUpToDate() {
		if err := os.RemoveAll(g.getDownloadPath()); err != nil {
			return "", fmt.Errorf("RemoveAll: %s", err.Error())
		}
		_, err := git.PlainClone(g.getDownloadPath(), false, &git.CloneOptions{
			URL:      g.Cfg.Git.Url,
			Progress: os.Stdout,
			Depth:    1,
		})
		if err != nil {
			return "", fmt.Errorf("PlainClone: %s", err.Error())
		}
		err = ioutil.WriteFile(g.getCacheFilePath(), []byte(strconv.Itoa(int(time.Now().Unix()))), 0644)
		if err != nil {
			return "", fmt.Errorf("WriteFile: %s", err.Error())
		}
	}

	refName = strings.TrimRight(refName, "/")
	content, err := g.getFileContent(refName, filePath)
	if err != nil {
		return "", fmt.Errorf("GitShow (%s): %s", g.getShowRef(refName, filePath), err.Error())
	}
	return g.renderContent(refType, content, refName, filePath)
}

func (g *GitHandler) getShowRef(branchName string, filePath string) string {
	dirPath := path.Join(".", g.Cfg.Git.WorkDir, filePath)
	if dirPath == "." {
		dirPath += "/"
	}
	return branchName + ":" + dirPath
}
