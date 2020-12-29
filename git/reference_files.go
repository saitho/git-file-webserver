package git

import (
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/godo.v2/glob"
)

func IsErrGitFileNotFound(err error) bool {
	return err.Error() == "exit status 128"
}

func (ref *Reference) getFileContent(filePath string) (string, error) {
	return ref.Client.runGitCommand(ref.Client.CurrentRepo, "show", ref.GetShowRef(filePath))
}

func (ref *Reference) isFolder(content string, filePath string) bool {
	contentLines := strings.Split(content, "\n")
	return contentLines[0] == "tree "+ref.GetShowRef(filePath)
}

func (ref *Reference) matchesFileGlob(pathWithParent string) (bool, error) {
	if len(ref.Client.Cfg.Files) == 0 {
		return true, nil
	}
	fileContent, err := ref.getFileContent(pathWithParent)
	if err != nil {
		return false, err
	}

	if ref.isFolder(fileContent, pathWithParent) {
		filesInFolder := strings.Split(fileContent, "\n")[2:]
		matchingFilesInFolder := ref.filterFiles(filesInFolder, pathWithParent)
		log.Tracef("-- path is folder - %s - %d\n", pathWithParent, len(matchingFilesInFolder))
		return len(matchingFilesInFolder) > 0, nil
	} else {
		for _, globString := range ref.Client.Cfg.Files {
			regExp := glob.Globexp(globString)
			matchPath := strings.TrimPrefix(pathWithParent, "/")
			match := regExp.Match([]byte(matchPath))
			log.Tracef("-- path is file - %s - %s- %v\n", matchPath, regExp, match)
			if match {
				return true, nil
			}
		}
	}
	return false, nil
}

func (ref *Reference) filterFiles(files []string, parentFolder string) []string {
	var newFiles []string
	for _, file := range files {
		file = strings.TrimPrefix(file, "/")
		if file == "" {
			continue
		}
		pathWithParent := path.Join(parentFolder, file)
		pathMatches, err := ref.matchesFileGlob(pathWithParent)
		if err != nil {
			panic("matchesFileGlob: " + err.Error())
		}

		if pathMatches {
			newFiles = append(newFiles, file)
		}
	}
	return newFiles
}
