//go:generate pkger

package main

import (
	"flag"
	"io"
	"net/http"
	"os"

	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"

	"github.com/saitho/static-git-file-server/config"
	"github.com/saitho/static-git-file-server/git"
	"github.com/saitho/static-git-file-server/webserver"
)

func initLoggers(cfg *config.Config) {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Panic(err)
	}

	var logLevel log.Level
	switch cfg.LogLevel {
	case "info":
		logLevel = log.InfoLevel
	case "warning":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	case "panic":
		logLevel = log.PanicLevel
	default:
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	formatter := new(log.TextFormatter)
	formatter.ForceColors = true
	log.SetFormatter(formatter)
}

func main() {
	_ = pkger.Include("/tmpl")

	port := flag.String("p", "80", "port to serve on")
	configPath := flag.String("c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}
	initLoggers(cfg)

	client := &git.Client{Cfg: cfg}
	server := webserver.Webserver{
		Port:       *port,
		ConfigPath: *configPath,
	}

	server.AddHandler(`^/(.*)/webhook/github`, webserver.GitHubWebHookEndpoint(client))
	if cfg.Display.Tags.VirtualTags.EnableSemverMajor {
		server.AddHandler(`^/(.*)/tag/(v?\d+)/-/(.*)$`, webserver.ResolveVirtualMajorTag(client))
		server.AddHandler(`^/(.*)/tag/(v?\d+)/?$`, webserver.ResolveVirtualMajorTag(client))
	}

	server.AddHandler(`^/(.*)/(branch|tag)/(.*)/-/(.*)$`, webserver.FileHandler(client))
	server.AddHandler(`^/(.*)/(branch|tag)/(.*)/?$`, webserver.FileHandler(client))

	// Redirect single-repo path to new multi-repo path
	server.AddHandler(`^/webhook/github`, func(resp *webserver.Response, req *webserver.Request) {
		http.Redirect(resp, req.Request, "/"+cfg.Git.Repositories[0].Slug+"/webhook/github", http.StatusPermanentRedirect)
	})
	server.AddHandler(`^/(branch|tag)/(.*)$`, func(resp *webserver.Response, req *webserver.Request) {
		http.Redirect(resp, req.Request, "/"+cfg.Git.Repositories[0].Slug+"/"+req.Params[0]+"/"+req.Params[1], http.StatusPermanentRedirect)
	})

	server.AddHandler(`^/(branch|tag)/?$`, func(resp *webserver.Response, req *webserver.Request) {
		http.Redirect(resp, req.Request, "/", http.StatusPermanentRedirect)
	})
	server.AddHandler(`^/$`, webserver.IndexHandler(client))
	server.Run()

}
