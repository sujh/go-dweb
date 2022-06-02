package dweb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	w          http.ResponseWriter
	r          *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string
}

type H map[string]any

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w:      w,
		r:      r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context) Data(status int, data []byte) {
	c.Status(status)
	c.w.Write(data)
}

func (c *Context) HTML(status int, content string) {
	c.Status(status)
	c.SetHeader("Content-Type", "text/html")
	c.w.Write([]byte(content))
}

func (c *Context) String(status int, format string, contents ...any) {
	c.Status(status)
	c.SetHeader("Content-Type", "text/plain")
	fmt.Fprintf(c.w, format, contents...)
}

func (c *Context) JSON(status int, obj any) {
	c.Status(status)
	c.SetHeader("Content-Type", "application/json")
	encoder := json.NewEncoder(c.w)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) Query(name string) string {
	return c.r.URL.Query().Get(name)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.w.WriteHeader(code)
}

func (c *Context) PostForm(key string) string {
	return c.r.FormValue(key)
}

func (c *Context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}
