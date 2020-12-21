//go:generate pkger

package main

import (
	"flag"
	"github.com/markbates/pkger"
	"github.com/saitho/static-git-file-server/rendering"
	"net/http"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/webserver"
)

func main() {
	pkger.Include("/tmpl")

	port := flag.String("p", "80", "port to serve on")
	configPath := flag.String("c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	gitHandler := git.GitHandler{Cfg: cfg}
	server := webserver.Webserver{
		Port:       *port,
		ConfigPath: *configPath,
	}

	handler := func(resp *webserver.Response, req *webserver.Request) {
		filePath := ""
		if len(req.Params) > 2 {
			filePath = req.Params[2]
		}
		content, err := gitHandler.ServePath(req.Params[0], req.Params[1], filePath)
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

	server.AddHandler(`^/(branch|tag)/(.*)/-/(.*)`, handler)
	server.AddHandler(`^/(branch|tag)/(.*)/?$`, handler)
	server.AddHandler(`^/$`, func(resp *webserver.Response, req *webserver.Request) {
		type IndexTmplParams struct {
			Cfg          *config.Config
			ShowBranches bool
			ShowTags     bool
			Branches     []string
			Tags         []git.GitTag
		}

		content, err := rendering.RenderTemplate("/tmpl/index.html", IndexTmplParams{
			Cfg:          cfg,
			ShowBranches: cfg.Display.Index.ShowBranches,
			ShowTags:     cfg.Display.Index.ShowTags,
			Branches:     gitHandler.GetBranches(),
			Tags:         gitHandler.GetTags(),
		})
		if err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}

		resp.HTML(http.StatusOK, content)
	})
	server.Run()

}
