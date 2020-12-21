package main

import (
	"bytes"
	"flag"
	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/webserver"
	"html/template"
	"net/http"
	"strings"
)

func main() {
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
		resp.HTML(http.StatusOK, content)
	}

	server.AddHandler(`^/(branch|tag)/(.*)/-/(.*)`, handler)
	server.AddHandler(`^/(branch|tag)/(.*)/?$`, handler)
	server.AddHandler(`^/$`, func(resp *webserver.Response, req *webserver.Request) {
		tplFuncMap := make(template.FuncMap)
		tplFuncMap["Split"] = strings.Split
		t, err := template.New("index.html").Funcs(tplFuncMap).ParseFiles("tmpl/index.html")
		if err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}
		var tpl bytes.Buffer
		type IndexTmplParams struct {
			Cfg          *config.Config
			ShowBranches bool
			ShowTags     bool
			Branches     []string
			Tags         []git.GitTag
		}
		params := IndexTmplParams{
			Cfg:          cfg,
			ShowBranches: cfg.Display.Index.ShowBranches,
			ShowTags:     cfg.Display.Index.ShowTags,
			Branches:     gitHandler.GetBranches(),
			Tags:         gitHandler.GetTags(),
		}

		if err := t.Execute(&tpl, params); err != nil {
			resp.Text(http.StatusInternalServerError, err.Error())
			return
		}

		resp.HTML(http.StatusOK, tpl.String())
	})
	server.Run()

}
