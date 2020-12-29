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
	"github.com/saitho/static-git-file-server/utils"
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
		// params: SLUG/tag/MAJORVERSION/-/PATH...
		var repoSlug, majorVersion, path string
		utils.Unpack(req.Params, &repoSlug, &majorVersion, nil, &path)

		if err := initRepo(client, repoSlug); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
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
	// params: SLUG/TYPE/REFNAME/-/PATH
	return func(resp *Response, req *Request) {
		var slug, refType, refName, filePath string
		utils.Unpack(req.Params, &slug, &refType, &refName, nil, &filePath)

		if err := initRepo(client, slug); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}

		reference := git.Reference{
			Client:   client,
			Type:     refType,
			Name:     strings.Trim(refName, "/"),
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
