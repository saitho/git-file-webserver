package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
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
	r.Header().Set("Content-Type", "text/plain")
	r.WriteHeader(code)

	io.WriteString(r, fmt.Sprintf("%s\n", body))
}
func (r *Response) HTML(code int, body string) {
	r.Header().Set("Content-Type", "text/html")
	r.WriteHeader(code)

	io.WriteString(r, fmt.Sprintf("%s\n", body))
}
