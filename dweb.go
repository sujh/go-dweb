package dweb

import (
	"log"
	"net/http"
)

type engine struct {
	router router
}

type HandlerFunc func(c *Context)

func New() *engine {
	return &engine{router: *newRouter()}
}

func (e *engine) GET(path string, handler HandlerFunc) {
	e.router.add("GET", path, handler)
}

func (e *engine) POST(path string, handler HandlerFunc) {
	e.router.add("POST", path, handler)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}

func (e *engine) Run(port string) {
	log.Fatal(http.ListenAndServe(port, e))
}
