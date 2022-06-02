package dweb

import "strings"

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChildren(part string) []*node {
	rst := make([]*node, 0)
	for _, c := range n.children {
		if c.part == part || c.isWild {
			rst = append(rst, c)
		}
	}
	return rst
}

func (n *node) matchChild(part string) *node {
	for _, c := range n.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

func (n *node) insert(pattern string, parts []string, depth int) {
	if depth == len(parts) {
		n.pattern = pattern
		return
	}
	part := parts[depth]
	c := n.matchChild(part)
	if c == nil {
		c = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, c)
	}
	c.insert(pattern, parts, depth+1)
}

func (n *node) search(parts []string, depth int) *node {
	if strings.HasPrefix(n.part, "*") {
		return n
	}
	if depth == len(parts) {
		if n.pattern != "" {
			return n
		}
		return nil
	}
	children := n.matchChildren(parts[depth])
	for _, c := range children {
		nc := c.search(parts, depth+1)
		if nc != nil {
			return nc
		}
	}
	return nil
}
