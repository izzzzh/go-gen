package tree

import (
	"errors"
	"strings"
)

var (
	errRouterNotFound = errors.New("router not found")
)

type (
	Node struct {
		pattern  []string
		part     string  // 路由中的一部分，例如 :lang
		children []*Node // 子节点，例如 [doc, tutorial, intro]
		isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
		item     any
	}

	Tree struct {
		root *Node
	}

	Result struct {
		Item   any
		Params map[string]string
	}
)

func NeeTree() *Tree {
	return &Tree{
		root: &Node{},
	}
}

func (t *Tree) Add(parts []string, handler any) error {
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
		if part[0] == ':' {
			params[part[1:]] = parts[i]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(parts[i:], "/")
			break
		}
	}
	ret.Item = item.item
	ret.Params = params
	return ret, nil
}

func (n *Node) add(parts []string, height int, handler any) error {
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
