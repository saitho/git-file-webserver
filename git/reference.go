package git

import (
	"path"

	"github.com/saitho/static-git-file-server/config"
)

type Reference struct {
	Client   *Client
	Cfg      config.Config
	Type     string
	Name     string
	FilePath string
}

func (ref *Reference) GetShowRef(filePath string) string {
	dirPath := path.Join(".", ref.Client.Cfg.Git.WorkDir, filePath)
	if dirPath == "." {
		dirPath += "/"
	}
	if ref.Type == "tag" {
		return "refs/tags/" + ref.Name + ":" + dirPath
	}
	return "origin/" + ref.Name + ":" + dirPath
}
