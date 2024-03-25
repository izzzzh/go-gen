package router

import (
	"errors"
	"github.com/valyala/fasthttp"
	"go-gen/gee/tree"
	"net/http"
	"path"
	"strings"
)

var (
	ErrInvalidMethod = errors.New("not a valid http method")
)

type (
	Route struct {
		Method  string
		Path    string
		Handler fasthttp.RequestHandler
	}
	Router struct {
		routers []*Route
		roots   map[string]*tree.Tree
	}
)

func New() *Router {
	return &Router{
		roots: map[string]*tree.Tree{},
	}
}

func (r *Router) Handle(ctx *fasthttp.RequestCtx) {
	reqPath := string(ctx.Path())
	method := string(ctx.Method())
	result, err := r.getRoute(method, reqPath)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}
	handler := result.Item.(fasthttp.RequestHandler)
	for k, v := range result.Params {
		ctx.SetUserValue(k, v)
	}
	handler(ctx)
}

func (r *Router) getRoute(method string, reqPath string) (tree.Result, error) {
	var ret tree.Result
	if !validMethod(method) {
		return ret, ErrInvalidMethod
	}
	parts := parseRouter(reqPath)
	searchTree := r.roots[method]

	search, err := searchTree.Search(parts)
	if err != nil {
		return tree.Result{}, err
	}
	return search, nil
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
		searchTree, ok := r.roots[method]
		if !ok {
			searchTree = tree.NeeTree()
		}

		parts := parseRouter(pattern)

		r.roots[method] = searchTree

		err := searchTree.Add(parts, route.Handler)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Router) AddRoute(router *Route) {
	r.routers = append(r.routers, router)
}
