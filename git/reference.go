package git

import (
	"path"

	"github.com/saitho/static-git-file-server/config"
)

type Reference struct {
	Client         *Client
	Cfg            config.Config
	Type           string
	Name           string
	FilePath       string
	FromVirtualTag string
}

func (ref *Reference) GetPath() string {
	if len(ref.FromVirtualTag) > 0 {
		return path.Join(ref.Client.CurrentRepo.Slug, "tag", ref.FromVirtualTag)
	}
	return path.Join(ref.Client.CurrentRepo.Slug, ref.Type, ref.Name)
}

func (ref *Reference) GetName() string {
	if len(ref.FromVirtualTag) > 0 {
		return ref.FromVirtualTag
	}
	return ref.Name
}

func (ref *Reference) GetShowRef(filePath string) string {
	dirPath := path.Join(".", ref.Client.CurrentRepo.WorkDir, filePath)
	if dirPath == "." {
		dirPath += "/"
	}
	if ref.Type == "tag" {
		return "refs/tags/" + ref.Name + ":" + dirPath
	}
	return "origin/" + ref.Name + ":" + dirPath
}
