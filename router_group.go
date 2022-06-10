package dweb

import (
	"net/http"
	"path"
)

type routerGroup struct {
	prefix      string
	middlewares []HandlerFunc
	router      *router
}

func GenTopRouterGroup() *routerGroup {
	return &routerGroup{prefix: "", router: newRouter()}
}

func (g *routerGroup) GET(pattern string, handle HandlerFunc) {
	g.addRoute("GET", pattern, handle)
}

func (g *routerGroup) POST(pattern string, handle HandlerFunc) {
	g.addRoute("POST", pattern, handle)
}

func (g *routerGroup) addRoute(method string, pattern string, handle HandlerFunc) {
	g.router.add(method, g.prefix+pattern, handle)
}

func (g *routerGroup) Group(prefix string) *routerGroup {
	cg := routerGroup{
		prefix: g.prefix + prefix,
		router: g.router,
	}
	return &cg
}

func (g *routerGroup) Use(mid HandlerFunc) {
	g.middlewares = append(g.middlewares, mid)
}

func (g *routerGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileHandler := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		if _, e := fs.Open(file); e != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileHandler.ServeHTTP(c.w, c.r)
	}
}

func (g *routerGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}
