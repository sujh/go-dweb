package dweb

import (
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

func (r *router) add(method string, pattern string, handler HandlerFunc) {
	r.handlers[routeKey(method, pattern)] = handler
	n := r.roots[method]
	if n == nil {
		n = &node{}
		r.roots[method] = n
	}
	parts := parsePattern(pattern)
	n.insert(pattern, parts, 0)
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := routeKey(c.Method, node.pattern)
		if handler, ok := r.handlers[key]; ok {
			handler(c)
			return
		}
	}
	c.String(http.StatusNotFound, "404 not found: %s for %s\n", c.Method, c.Path)
}

func parsePattern(pattern string) []string {
	rawParts := strings.Split(pattern, "/")
	rst := make([]string, 0)
	for _, p := range rawParts {
		if p != "" {
			rst = append(rst, p)
			if p[0] == '*' {
				break
			}
		}
	}
	return rst
}

func (r *router) getRoute(method string, path string) (n *node, params map[string]string) {
	root := r.roots[method]
	params = make(map[string]string, 0)
	if root != nil {
		n = root.search(parsePattern(path), 0)
		if n != nil {
			np := parsePattern(n.pattern)
			pp := strings.FieldsFunc(path, func(c rune) bool {
				return c == '/'
			})
			for idx, part := range np {
				if part[0] == ':' {
					params[part[1:]] = pp[idx]
				} else if part[0] == '*' && len(part) > 1 {
					params[part[1:]] = strings.Join(pp[idx:], "/")
					break
				}
			}
			return n, params
		}
	}
	return nil, nil
}

func routeKey(method string, path string) string {
	return method + "-" + path
}
