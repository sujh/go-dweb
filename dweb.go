package dweb

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type engine struct {
	groups         []*routerGroup
	topRouterGroup *routerGroup
}

type HandlerFunc func(c *Context)

func New() *engine {
	tg := GenTopRouterGroup()
	e := engine{topRouterGroup: tg}
	e.groups = append(e.groups, tg)
	return &e
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	for _, g := range e.groups {
		if strings.HasPrefix(r.URL.Path, g.prefix) {
			c.handlers = append(c.handlers, g.middlewares...)
		}
	}
	e.topRouterGroup.router.handle(c)
}

func (e *engine) Run(port string) {
	log.Fatal(http.ListenAndServe(port, e))
}

func (e *engine) TopRouterGroup() *routerGroup {
	return e.topRouterGroup
}

func (e *engine) GenerateAndRecordRouterGroup(parent *routerGroup, prefix string) *routerGroup {
	newGroup := parent.Group(prefix)
	e.groups = append(e.groups, newGroup)
	return newGroup
}

// Default middleware
func Logger(c *Context) {
	t := time.Now()
	c.Yield()
	log.Printf("[%d] %s in %v", c.StatusCode, c.r.RequestURI, time.Since(t))
}
