package git

import (
	"bytes"
	"github.com/saitho/static-git-file-server/config"
	"html/template"
	"path"
	"strings"
)

func (g *GitHandler) renderContent(refType string, content string, refName string, filePath string) (string, error) {
	if !g.isFolder(content, refName, filePath) {
		return content, nil
	}
	// Render HTML template for folders and list all files there

	contentLines := strings.Split(content, "\n")
	t, err := template.ParseFiles("tmpl/dir.html") // Parse template file.
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	var pathChunks = strings.Split(filePath, "/")
	parentPath := ""
	if len(pathChunks) > 0 {
		parentPath = strings.Join(pathChunks[:len(pathChunks)-1], "/")
	}
	var data = TmplParams{
		Cfg:            g.Cfg,
		RefType:        refType,
		Path:           filePath,
		ParentPath:     parentPath,
		FullPath:       path.Join(g.Cfg.Git.WorkDir, filePath),
		FullParentPath: path.Join(g.Cfg.Git.WorkDir, parentPath),
		Ref:            refName,
		Files:          g.filterFiles(contentLines[2:], filePath, refName),
	}
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}

type TmplParams struct {
	Cfg            config.Config
	RefType        string
	Path           string
	ParentPath     string
	FullPath       string
	FullParentPath string
	Ref            string
	Files          []string
}
