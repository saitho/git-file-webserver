package webserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/rendering"
)

func initRepo(client *git.Client, repoSlug string) error {
	repoList := client.Cfg.Git.Repositories
	if len(repoSlug) > 0 {
		_ = client.SelectRepository(repoSlug)
		repoList = []*config.RepoConfig{client.CurrentRepo}
	}
	// Update all repositories
	for _, repo := range repoList {
		if !repo.IsUpToDate() {
			log.Debugf("Downloading repository as repo was not cloned yet or is outdated by cache time.")
			if err := client.DownloadRepository(repo); err != nil {
				return err
			}
		}
	}
	return nil
}

func ResolveVirtualMajorTag(client *git.Client) func(resp *Response, req *Request) {
	return func(resp *Response, req *Request) {
		repoSlug := req.Params[0]
		if err := initRepo(client, repoSlug); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		majorVersion := req.Params[1]
		path := ""
		if len(req.Params) > 2 {
			path = req.Params[2]
		}

		latestTag, err := git.ResolveVirtualTag(client, majorVersion)
		if err != nil {
			log.Errorf("Unable to resolve tag %s", majorVersion)
			resp.Text(http.StatusInternalServerError, fmt.Sprintf("Unable to resolve tag %s", majorVersion))
			return
		}

		req.Params = []string{repoSlug, "tag", latestTag.Tag, path}
		FileHandler(client)(resp, req)
	}
}

func FileHandler(client *git.Client) func(resp *Response, req *Request) {
	fmt.Println("FileHandler")
	return func(resp *Response, req *Request) {
		if err := initRepo(client, req.Params[0]); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		filePath := ""
		if len(req.Params) > 3 {
			filePath = req.Params[3]
		}

		reference := git.Reference{
			Client:   client,
			Type:     req.Params[1],
			Name:     strings.Trim(req.Params[2], "/"),
			FilePath: strings.Trim(filePath, "/"),
		}
		content, err := reference.Render()
		if err != nil {
			if git.IsErrGitFileNotFound(err) {
				log.Warningf("File '%s' not found on %s %s (repo: %s).", reference.FilePath, reference.Type, reference.Name, client.CurrentRepo.Slug)
				resp.Text(http.StatusNotFound, "Requested file not found.")
				return
			}
			log.Errorf("Unexpected error during Git file lookup: %s", err)
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		resp.Auto(http.StatusOK, content)
	}
}

func IndexHandler(client *git.Client) func(resp *Response, req *Request) {
	return func(resp *Response, req *Request) {
		if err := initRepo(client, ""); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		type IndexTmplParams struct {
			Cfg          *config.Config
			ShowBranches bool
			ShowTags     bool
			Branches     []string
			Tags         []git.GitTag
			LastUpdate   time.Time
			Client       *git.Client
		}

		cfg := client.Cfg
		content, err := rendering.RenderTemplate("/tmpl/index.html", IndexTmplParams{
			Cfg:          cfg,
			ShowBranches: cfg.Display.Index.ShowBranches,
			ShowTags:     cfg.Display.Index.ShowTags,
			Client:       client,
		})
		if err != nil {
			log.Errorf("Unable to render index template: %s", err)
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}

		resp.HTML(http.StatusOK, content)
	}
}
