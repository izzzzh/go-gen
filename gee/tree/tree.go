package tree

import (
	"errors"
	"github.com/valyala/fasthttp"
	"strings"
)

var (
	errRouterNotFound = errors.New("router not found")
)

type (
	Node struct {
		pattern  []string
		part     string
		children []*Node
		isWild   bool
		item     fasthttp.RequestHandler
	}

	Tree struct {
		root *Node
	}

	Result struct {
		Item   fasthttp.RequestHandler
		Params map[string]string
	}
)

func NewTree() *Tree {
	return &Tree{
		root: &Node{},
	}
}

func (t *Tree) Add(parts []string, handler fasthttp.RequestHandler) error {
	node := t.root
	err := node.add(parts, 0, handler)
	return err
}

func (t *Tree) Search(parts []string) (Result, error) {
	var ret Result
	node := t.root
	item := node.search(parts, 0)
	if item == nil {
		return ret, errRouterNotFound
	}
	pattern := item.pattern
	params := make(map[string]string)
	for i := range pattern {
		part := pattern[i]
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			params[part[1:]] = parts[i]
		}
		if part[0] == '*' && len(part) > 1 {
			var builder strings.Builder
			builder.Grow(len(parts[i:]))
			for j := i; j < len(parts); j++ {
				builder.WriteString(parts[j])
				if j != len(parts)-1 {
					builder.WriteString("/")
				}
			}
			params[part[1:]] = builder.String()
			break
		}
	}
	ret.Item = item.item
	ret.Params = params
	return ret, nil
}

func (n *Node) add(parts []string, height int, handler fasthttp.RequestHandler) error {
	if len(parts) == height {
		n.item = handler
		n.pattern = parts
		return nil
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{part: part, isWild: part[0] == ':' || part[0] == '*', children: nil}
		n.children = append(n.children, child)
	}
	return child.add(parts, height+1, handler)
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.item == nil {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
