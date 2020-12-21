package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	glob "gopkg.in/godo.v2/glob"
)

const DownloadLocation = "./git_downloads"

type GitHandler struct {
	Cfg Config
}

func (g *GitHandler) getCacheFilePath() string {
	return path.Join(g.getDownloadPath() + ".cache")
}
func (g *GitHandler) getDownloadPath() string {
	return path.Join(DownloadLocation, g.Cfg.Git.Url)
}

func (g *GitHandler) isUpToDate() bool {
	// File does not exist
	_, err := os.Stat(g.getDownloadPath())
	if os.IsNotExist(err) {
		return false
	}

	// Check download date
	cacheFile, _ := ioutil.ReadFile(g.getCacheFilePath())
	cacheTime, _ := strconv.Atoi(string(cacheFile))
	if time.Now().Unix() > int64(cacheTime+g.Cfg.Git.CacheTime) {
		return false
	}
	return true
}

func (g *GitHandler) ServePath() Handler {
	return func(resp *Response, req *Request) {
		if !g.isUpToDate() {
			if err := os.RemoveAll(g.getDownloadPath()); err != nil {
				panic("RemoveAll: " + err.Error())
			}
			_, err := git.PlainClone(g.getDownloadPath(), false, &git.CloneOptions{
				URL:      g.Cfg.Git.Url,
				Progress: os.Stdout,
				Depth:    1,
			})
			if err != nil {
				panic("PlainClone: " + err.Error())
			}
			err = ioutil.WriteFile(g.getCacheFilePath(), []byte(strconv.Itoa(int(time.Now().Unix()))), 0644)
			if err != nil {
				panic("WriteFile: " + err.Error())
			}
		}

		filePath := ""
		if len(req.Params) >= 3 {
			filePath = strings.TrimRight(req.Params[2], "/")
		}

		branchName := strings.TrimRight(req.Params[1], "/")
		content, err := g.getFileContent(branchName, filePath)
		if err != nil {
			resp.Text(500, "GitShow (" + g.getShowRef(branchName, filePath) + "): " + err.Error())
			return
		}
		g.parseContent(req.Params[0], content, resp, branchName, filePath, branchName)
	}
}

func (g *GitHandler) getFileContent(branchName string, filePath string) (string, error) {
	cmd:= exec.Command("git", "show", g.getShowRef(branchName, filePath))
	cmd.Dir = g.getDownloadPath()
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type TmplParams struct {
	Cfg Config
	RefType string
	Path string
	ParentPath string
	FullPath string
	FullParentPath string
	Ref string
	Files []string
}

func (g *GitHandler) getShowRef(branchName string, filePath string) string {
	dirPath := path.Join(".", g.Cfg.Git.WorkDir, filePath)
	if dirPath == "." {
		dirPath += "/"
	}
	return branchName + ":" + dirPath
}

func (g *GitHandler) isFolder(content string, refName string, filePath string) bool {
	contentLines := strings.Split(content, "\n")
	return contentLines[0] == "tree "  + g.getShowRef(refName, filePath)
}

func (g *GitHandler) parseContent(refType string, content string, resp *Response, refName string, filePath string, branchName string) {
	contentLines := strings.Split(content, "\n")
	if g.isFolder(content, refName, filePath) {
		// Evaluate as folder
		t, err := template.ParseFiles("tmpl/dir.html")  // Parse template file.
		if err != nil {
			resp.Text(500, "ParseFiles: " + err.Error())
			return
		}

		var tpl bytes.Buffer
		var pathChunks = strings.Split(filePath, "/")
		parentPath := ""
		if len(pathChunks) > 0 {
			parentPath = strings.Join(pathChunks[:len(pathChunks)-1], "/")
		}
		var data = TmplParams{
			Cfg: g.Cfg,
			RefType: refType,
			Path: filePath,
			ParentPath: parentPath,
			FullPath: path.Join(g.Cfg.Git.WorkDir, filePath),
			FullParentPath: path.Join(g.Cfg.Git.WorkDir, parentPath),
			Ref: refName,
			Files: g.filterFiles(contentLines[2:], filePath, branchName),
		}
		if err := t.Execute(&tpl, data); err != nil {
			resp.Text(500, "Execute: " + err.Error())
			return
		}

		resp.HTML(200, tpl.String())
		return
	}
	resp.Text(200, content)
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
