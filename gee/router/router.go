package router

import (
	"errors"
	"github.com/valyala/fasthttp"
	"net/http"
	"path"
	"strings"
)

type Node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*Node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

type Route struct {
	Method  string
	Path    string
	Handler fasthttp.RequestHandler
}

type Router struct {
	routers  []*Route
	roots    map[string]*Node
	handlers map[string]fasthttp.RequestHandler
}

func New() *Router {
	return &Router{
		roots:    map[string]*Node{},
		routers:  []*Route{},
		handlers: map[string]fasthttp.RequestHandler{},
	}
}

func (n *Node) matchChild(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{part: part, isWild: part[0] == ':' || part[0] == '*', children: nil}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
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

func (r *Router) Handle(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	method := string(ctx.Method())
	node, params := r.getRoute(method, path)
	if node != nil {
		key := method + "-" + node.pattern
		handler := r.handlers[key]
		for k, v := range params {
			ctx.SetUserValue(k, v)
		}
		handler(ctx)
	} else {
		ctx.NotFound()
	}
	return
}

func (r *Router) getRoute(method string, path string) (*Node, map[any]any) {
	searchParts := parseRouter(path)
	params := make(map[any]any)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parseRouter(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func parseRouter(pattern string) []string {
	pattern = path.Clean(pattern)
	pathArr := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range pathArr {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func validMethod(method string) bool {
	return method == http.MethodDelete || method == http.MethodGet ||
		method == http.MethodHead || method == http.MethodOptions ||
		method == http.MethodPatch || method == http.MethodPost ||
		method == http.MethodPut
}

func (r *Router) BindRoute() error {
	for _, route := range r.routers {
		method := route.Method
		pattern := route.Path
		if !validMethod(method) {
			return errors.New("method not allow")
		}
		parts := parseRouter(route.Path)
		key := route.Method + "-" + route.Path
		_, ok := r.roots[route.Method]
		if !ok {
			r.roots[method] = &Node{}
		}
		r.roots[method].insert(pattern, parts, 0)
		r.handlers[key] = route.Handler
	}

	return nil
}

func (r *Router) AddRoute(router *Route) {
	r.routers = append(r.routers, router)
}
