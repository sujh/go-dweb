package dweb

import (
	"fmt"
	"log"
	"net/http"
)

type engine struct {
	router map[string]HandlerFunc
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func New() *engine {
	return &engine{router: make(map[string]HandlerFunc)}
}

func (e *engine) addRoute(verb string, path string, handler HandlerFunc) {
	e.router[verb+"-"+path] = handler
}

func (e *engine) GET(path string, handler HandlerFunc) {
	e.addRoute("GET", path, handler)
}

func (e *engine) POST(path string, handler HandlerFunc) {
	e.addRoute("POST", path, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := e.router[r.Method+"-"+r.URL.Path]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "NOT FOUND: %s", r.URL)
	}
}

func (e *engine) Run(port string) {
	log.Fatal(http.ListenAndServe(port, e))
}
