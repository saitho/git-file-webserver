package git

import (
	"os/exec"
	"path"
	"strings"

	"gopkg.in/godo.v2/glob"
)

func IsErrGitFileNotFound(err error) bool {
	return err.Error() == "exit status 128"
}

func (g *GitHandler) getFileContent(branchName string, filePath string) (string, error) {
	cmd := exec.Command("git", "show", g.getShowRef(branchName, filePath))
	cmd.Dir = g.getDownloadPath()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (g *GitHandler) isFolder(content string, refName string, filePath string) bool {
	contentLines := strings.Split(content, "\n")
	return contentLines[0] == "tree "+g.getShowRef(refName, filePath)
}

func (g *GitHandler) matchesFileGlob(pathWithParent string, branchName string) (bool, error) {
	fileContent, err := g.getFileContent(branchName, pathWithParent)
	if err != nil {
		return false, err
	}

	if g.isFolder(fileContent, branchName, pathWithParent) {
		filesInFolder := strings.Split(fileContent, "\n")[2:]
		matchingFilesInFolder := g.filterFiles(filesInFolder, pathWithParent, branchName)
		// fmt.Printf("-- path is folder - %s - %d\n", pathWithParent, len(matchingFilesInFolder))
		return len(matchingFilesInFolder) > 0, nil
	} else {
		for _, globString := range g.Cfg.Files {
			regExp := glob.Globexp(globString)
			matchPath := strings.TrimLeft(pathWithParent, "/")
			match := regExp.Match([]byte(matchPath))
			// fmt.Printf("-- path is file - %s - %s- %b\n", matchPath, regExp, match)
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func (g *GitHandler) filterFiles(files []string, parentFolder string, branchName string) []string {
	var newFiles []string
	for _, file := range files {
		if file == "" {
			continue
		}
		pathWithParent := path.Join(parentFolder, file)
		pathMatches, err := g.matchesFileGlob(pathWithParent, branchName)
		if err != nil {
			panic("matchesFileGlob: " + err.Error())
		}

		if pathMatches {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}
