//go:generate pkger

package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/markbates/pkger"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/rendering"
	"github.com/saitho/static-git-file-server/webserver"
)

func main() {
	_ = pkger.Include("/tmpl")

	port := flag.String("p", "80", "port to serve on")
	configPath := flag.String("c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	client := &git.Client{Cfg: cfg}
	initRepo := func() error {
		if !client.IsUpToDate() {
			if err := client.DownloadRepository(); err != nil {
				return err
			}
		}
		return nil
	}
	server := webserver.Webserver{
		Port:       *port,
		ConfigPath: *configPath,
	}

	handler := func(resp *webserver.Response, req *webserver.Request) {
		if err := initRepo(); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		filePath := ""
		if len(req.Params) > 2 {
			filePath = req.Params[2]
		}

		reference := git.Reference{
			Client:   client,
			Type:     req.Params[0],
			Name:     strings.Trim(req.Params[1], "/"),
			FilePath: strings.Trim(filePath, "/"),
		}
		content, err := reference.Render()
		if err != nil {
			if git.IsErrGitFileNotFound(err) {
				resp.Text(http.StatusNotFound, "Requested file not found.")
				return
			}
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		resp.Auto(http.StatusOK, content)
	}

	resolveVirtualMajorTag := func(resp *webserver.Response, req *webserver.Request) {
		majorVersion := req.Params[0]
		path := ""
		if len(req.Params) > 1 {
			path = req.Params[1]
		}

		latestTag, err := git.ResolveVirtualTag(client, majorVersion)
		if err != nil {
			resp.Text(http.StatusInternalServerError, fmt.Sprintf("Unable to resolve tag %s", majorVersion))
			return
		}

		req.Params = []string{"tag", latestTag.Tag, path}
		handler(resp, req)
	}

	server.AddHandler(`^/webhook/github`, webserver.GitHubWebHookEndpoint(cfg, client))
	if cfg.Display.Tags.VirtualTags.EnableSemverMajor {
		server.AddHandler(`^/tag/(v?\d+)/-/(.*)`, resolveVirtualMajorTag)
		server.AddHandler(`^/tag/(v?\d+)/?$`, resolveVirtualMajorTag)
	}

	server.AddHandler(`^/(branch|tag)/(.*)/-/(.*)`, handler)
	server.AddHandler(`^/(branch|tag)/(.*)/?$`, handler)
	server.AddHandler(`^/$`, func(resp *webserver.Response, req *webserver.Request) {
		if err := initRepo(); err != nil {
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
		}

		tags := client.GetTags()
		if strings.ToLower(cfg.Display.Tags.Order) == "asc" {
			// Reverse array
			for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}

		if cfg.Display.Tags.VirtualTags.EnableSemverMajor {
			tags = git.InsertVirtualTags(tags)
		}

		content, err := rendering.RenderTemplate("/tmpl/index.html", IndexTmplParams{
			Cfg:          cfg,
			ShowBranches: cfg.Display.Index.ShowBranches,
			ShowTags:     cfg.Display.Index.ShowTags,
			Branches:     client.GetBranches(),
			Tags:         tags,
			LastUpdate:   time.Unix(client.GetUpdatedTime(), 0),
		})
		if err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}

		resp.HTML(http.StatusOK, content)
	})
	server.Run()

}
