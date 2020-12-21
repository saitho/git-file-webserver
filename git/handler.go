package git

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/saitho/static-git-file-server/config"
)

const DownloadLocation = "./git_downloads"

type GitHandler struct {
	Cfg *config.Config
}

func (g *GitHandler) getCacheFilePath() string {
	return path.Join(g.getDownloadPath() + ".cache")
}
func (g *GitHandler) getDownloadPath() string {
	return path.Join(DownloadLocation, g.Cfg.Git.Url)
}

func (g *GitHandler) runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = g.getDownloadPath()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (g *GitHandler) ServePath(refType string, refName string, filePath string) (string, error) {
	refName = strings.TrimSuffix(refName, "/")
	content, err := g.getFileContent(refName, filePath)
	if err != nil {
		if IsErrGitFileNotFound(err) {
			return "", err
		}
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
