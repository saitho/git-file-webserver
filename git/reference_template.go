package git

import (
	"fmt"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/saitho/static-git-file-server/rendering"
)

func (ref *Reference) Render() (string, error) {
	content, err := ref.getFileContent(ref.FilePath)
	if err != nil {
		if IsErrGitFileNotFound(err) {
			return "", err
		}
		return "", fmt.Errorf("GitShow (%s): %s", ref.GetShowRef(ref.FilePath), err.Error())
	}

	if !ref.isFolder(content, ref.FilePath) {
		log.Debugf("Serving file content of '%s' from %s %s.", ref.FilePath, ref.Type, ref.Name)
		return content, nil
	}
	log.Debugf("Rendering contents of folder (path: '%s', %s: %s)", ref.FilePath, ref.Type, ref.Name)
	// Render HTML template for folders and list all files there

	contentLines := strings.Split(content, "\n")
	var pathChunks = strings.Split(ref.FilePath, "/")
	parentPath := ""
	if len(pathChunks) > 0 {
		parentPath = strings.Join(pathChunks[:len(pathChunks)-1], "/")
	}

	return rendering.RenderTemplate("/tmpl/dir.html", TmplParams{
		Ref:            ref,
		ParentPath:     parentPath,
		FullPath:       path.Join(ref.Client.Cfg.Git.WorkDir, ref.FilePath),
		FullParentPath: path.Join(ref.Client.Cfg.Git.WorkDir, parentPath),
		Files:          ref.filterFiles(contentLines[2:], ref.FilePath),
	})
}

type TmplParams struct {
	Ref            *Reference
	ParentPath     string
	FullPath       string
	FullParentPath string
	Files          []string
}
