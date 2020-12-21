package webserver

import (
	"log"
	"net/http"
)

type Webserver struct {
	Port       string
	ConfigPath string

	handlers map[string]Handler
}

func (w *Webserver) AddHandler(pattern string, handler Handler) {
	if w.handlers == nil {
		w.handlers = map[string]Handler{}
	}
	w.handlers[pattern] = handler
}

func (w *Webserver) Run() {
	app := NewRequestHandler()
	for pattern, handler := range w.handlers {
		app.Handle(pattern, handler)
	}

	log.Printf("Serving with config at %s on HTTP port: %s\n", w.ConfigPath, w.Port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+w.Port, app))
}
