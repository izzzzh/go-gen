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
		roots map[string]*tree.Tree
	}
	Middleware func(next fasthttp.RequestHandler) fasthttp.RequestHandler
)

func New() *Router {
	return &Router{
		roots: map[string]*tree.Tree{},
	}
}

// handleError 封装错误处理逻辑
func handleError(ctx *fasthttp.RequestCtx, err error, statusCode int) {
	ctx.Error(err.Error(), statusCode)
}

func (r *Router) Handle(ctx *fasthttp.RequestCtx) {
	reqPath := string(ctx.Path())
	method := string(ctx.Method())
	result, err := r.getRoute(method, reqPath)
	if err != nil {
		handleError(ctx, err, http.StatusInternalServerError)
		return
	}
	requestHandler := result.Item
	for k, v := range result.Params {
		ctx.SetUserValue(k, v)
	}
	requestHandler(ctx)
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

// parseRouter 优化路径解析逻辑
func parseRouter(pattern string) []string {
	pattern = path.Clean(pattern)
	if pattern == "/" {
		return []string{""} // 显式处理路径以 "/" 开头的情况
	}
	pathArr := strings.Split(pattern, "/")
	parts := make([]string, 0, len(pathArr))
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

func (r *Router) BindRoute(route *Route, middlewares []Middleware) error {
	method := route.Method
	pattern := route.Path
	if !validMethod(method) {
		return errors.New("method not allow")
	}
	searchTree, ok := r.roots[method]
	if !ok {
		searchTree = tree.NewTree()
	}
	parts := parseRouter(pattern)

	handle := r.WrappedHandle(route, middlewares)

	r.roots[method] = searchTree

	err := searchTree.Add(parts, handle)
	if err != nil {
		return err
	}
	return nil
}

func (r *Router) WrappedHandle(router *Route, middlewares []Middleware) fasthttp.RequestHandler {
	handle := router.Handler
	for _, middleware := range middlewares {
		handle = middleware(handle)
	}
	return handle
}
