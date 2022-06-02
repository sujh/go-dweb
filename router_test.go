package dweb

import (
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.add("GET", "/", nil)
	r.add("GET", "/", nil)
	r.add("GET", "/hello/*name", nil)
	r.add("GET", "/hi/:name", nil)
	r.add("GET", "/hi/:name/:address", nil)
	r.add("GET", "/specific/path", nil)
	r.add("POST", "/any/*", nil)
	return r
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	cases := []struct {
		inputMethod string
		inputPath   string
		wantPattern string
		wantParams  map[string]string
	}{
		{"GET", "/", "/", map[string]string{}},
		{"GET", "/hello/spike/bbb", "/hello/*name", map[string]string{"name": "spike/bbb"}},
		{"GET", "/hi/spike", "/hi/:name", map[string]string{"name": "spike"}},
		{"GET", "/hi/spike/NewYork", "/hi/:name/:address", map[string]string{"name": "spike", "address": "NewYork"}},
		{"GET", "/specific/path", "/specific/path", map[string]string{}},
		{"POST", "/any/x/y/z", "/any/*", map[string]string{}},
	}
	for _, c := range cases {
		node, matching := r.getRoute(c.inputMethod, c.inputPath)
		if node.pattern != c.wantPattern {
			t.Errorf("Input: %s %s\n matched pattern should be %s, but got %s", c.inputMethod, c.inputPath, c.wantPattern, node.pattern)
		}
		if !reflect.DeepEqual(matching, c.wantParams) {
			t.Errorf("Input: %s %s\n Params should be %s, but got %s", c.inputMethod, c.inputPath, c.wantParams, matching)
		}
	}
}
