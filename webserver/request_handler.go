package webserver

import (
	"github.com/saitho/static-git-file-server/config"
	"io"
	"net/http"
	"regexp"

	"github.com/gabriel-vasile/mimetype"
)

type Handler func(*Response, *Request)

type Route struct {
	Pattern *regexp.Regexp
	Handler Handler
}

type RequestHandler struct {
	Routes       []Route
	DefaultRoute Handler
}

func NewRequestHandler() *RequestHandler {
	app := &RequestHandler{
		DefaultRoute: func(resp *Response, req *Request) {
			resp.Text(http.StatusNotFound, "Not found")
		},
	}

	return app
}

func (a *RequestHandler) Handle(pattern string, handler Handler) {
	re := regexp.MustCompile(pattern)
	route := Route{Pattern: re, Handler: handler}

	a.Routes = append(a.Routes, route)
}

func (a *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := &Request{Request: r}
	resp := &Response{w}

	for _, rt := range a.Routes {
		if matches := rt.Pattern.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
			if len(matches) > 1 {
				req.Params = matches[1:]
			}

			rt.Handler(resp, req)
			return
		}
	}

	a.DefaultRoute(resp, req)
}

type Request struct {
	*http.Request
	Params []string
}

type Response struct {
	http.ResponseWriter
}

func (r *Response) Text(code int, body string) {
	r.send(code, "text/plain", body)
}

func (r *Response) HTML(code int, body string) {
	r.send(code, "text/html", body)
}

func (r *Response) Auto(code int, body string) {
	contentType := mimetype.Detect([]byte(body)).String()
	r.send(code, contentType, body)
}

func (r *Response) send(code int, contentType string, body string) {
	r.Header().Set("Content-Type", contentType)
	r.Header().Set("Git-File-Webserver-Version", config.VERSION)
	r.WriteHeader(code)

	_, _ = io.WriteString(r, body)
}
