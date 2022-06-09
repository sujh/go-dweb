package dweb

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
