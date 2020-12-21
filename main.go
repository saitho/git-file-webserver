package main

import (
	"flag"
	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/webserver"
	"net/http"
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
			resp.Text(500, err.Error())
			return
		}
		resp.HTML(200, content)
	}

	server.AddHandler(`^/(branch|tag)/(.*)/-/(.*)`, handler)
	server.AddHandler(`^/(branch|tag)/(.*)/?$`, handler)
	server.AddHandler(`^/$`, func(resp *webserver.Response, req *webserver.Request) {
		resp.Text(http.StatusOK, "index")
	})
	server.Run()

}
