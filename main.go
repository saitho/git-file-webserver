package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("p", "80", "port to serve on")
	configPath := flag.String("c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	app := NewRequestHandler()
	git := GitHandler{Cfg: cfg}
	app.Handle(`^/(branch|tag)/(.*)/-/(.*)`, git.ServePath())
	app.Handle(`^/(branch|tag)/(.*)/?$`, git.ServePath())

	app.Handle(`^/$`, func(resp *Response, req *Request) {
		resp.Text(http.StatusOK, "index")
	})

	log.Printf("Serving with config at %s on HTTP port: %s\n", *configPath, *port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+ *port, app))

}
