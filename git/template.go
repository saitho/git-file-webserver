package git

import (
	"path"
	"strings"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/rendering"
)

func (g *GitHandler) renderContent(refType string, content string, refName string, filePath string) (string, error) {
	if !g.isFolder(content, refName, filePath) {
		return content, nil
	}
	// Render HTML template for folders and list all files there

	contentLines := strings.Split(content, "\n")
	var pathChunks = strings.Split(filePath, "/")
	parentPath := ""
	if len(pathChunks) > 0 {
		parentPath = strings.Join(pathChunks[:len(pathChunks)-1], "/")
	}

	return rendering.RenderTemplate("/tmpl/dir.html", TmplParams{
		Cfg:            g.Cfg,
		RefType:        refType,
		Path:           filePath,
		ParentPath:     parentPath,
		FullPath:       path.Join(g.Cfg.Git.WorkDir, filePath),
		FullParentPath: path.Join(g.Cfg.Git.WorkDir, parentPath),
		Ref:            refName,
		Files:          g.filterFiles(contentLines[2:], filePath, refName),
	})
}

type TmplParams struct {
	Cfg            *config.Config
	RefType        string
	Path           string
	ParentPath     string
	FullPath       string
	FullParentPath string
	Ref            string
	Files          []string
}
