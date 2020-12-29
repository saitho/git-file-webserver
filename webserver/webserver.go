package webserver

import (
	"net/http"
	"regexp"

	log "github.com/sirupsen/logrus"
)

type Webserver struct {
	Port       string
	ConfigPath string

	routes []*Route
}

func (w *Webserver) AddHandler(pattern string, handler Handler) {
	w.routes = append(w.routes, &Route{
		Pattern: regexp.MustCompile(pattern),
		Handler: handler,
	})
}

func (w *Webserver) Run() {
	app := NewRequestHandler()
	for _, route := range w.routes {
		app.Handle(route)
	}

	log.Infof("Serving with config at %s on HTTP port: %s\n", w.ConfigPath, w.Port)
	err := http.ListenAndServe("0.0.0.0:"+w.Port, app)
	if err != nil {
		log.Fatal(err)
	}
}
