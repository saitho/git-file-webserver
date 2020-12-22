package git

import (
	"path"
	"strings"

	"gopkg.in/godo.v2/glob"
)

func IsErrGitFileNotFound(err error) bool {
	return err.Error() == "exit status 128"
}

func (g *GitHandler) getFileContent(refType string, refName string, filePath string) (string, error) {
	return g.runGitCommand("show", g.getShowRef(refType, refName, filePath))
}

func (g *GitHandler) isFolder(content string, refType string, refName string, filePath string) bool {
	contentLines := strings.Split(content, "\n")
	return contentLines[0] == "tree "+g.getShowRef(refType, refName, filePath)
}

func (g *GitHandler) matchesFileGlob(pathWithParent string, refType string, refName string) (bool, error) {
	if len(g.Cfg.Files) == 0 {
		return true, nil
	}
	fileContent, err := g.getFileContent(refType, refName, pathWithParent)
	if err != nil {
		return false, err
	}

	if g.isFolder(fileContent, refType, refName, pathWithParent) {
		filesInFolder := strings.Split(fileContent, "\n")[2:]
		matchingFilesInFolder := g.filterFiles(filesInFolder, pathWithParent, refType, refName)
		// fmt.Printf("-- path is folder - %s - %d\n", pathWithParent, len(matchingFilesInFolder))
		return len(matchingFilesInFolder) > 0, nil
	} else {
		for _, globString := range g.Cfg.Files {
			regExp := glob.Globexp(globString)
			matchPath := strings.TrimPrefix(pathWithParent, "/")
			match := regExp.Match([]byte(matchPath))
			// fmt.Printf("-- path is file - %s - %s- %b\n", matchPath, regExp, match)
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func (g *GitHandler) filterFiles(files []string, parentFolder string, refType string, refName string) []string {
	var newFiles []string
	for _, file := range files {
		if file == "" {
			continue
		}
		pathWithParent := path.Join(parentFolder, file)
		pathMatches, err := g.matchesFileGlob(pathWithParent, refType, refName)
		if err != nil {
			panic("matchesFileGlob: " + err.Error())
		}

		if pathMatches {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}
