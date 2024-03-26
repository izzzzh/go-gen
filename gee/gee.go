package gee

import (
	"github.com/valyala/fasthttp"
	"go-gen/gee/handler"
	"go-gen/gee/router"
)

type engine struct {
	conf        string
	router      *router.Router
	routers     []*router.Route
	middlewares []router.Middleware
}

func NewEngine() *engine {
	return &engine{
		router:      router.New(),
		middlewares: []router.Middleware{},
	}
}

func (e *engine) Start() {

	e.middlewares = append(e.middlewares, handler.LoggingMiddleware)
	e.middlewares = append(e.middlewares, handler.ErrorMiddleware)

	err := e.BindRoutes()
	if err != nil {
		panic(err)
	}
	e.Run(":8080")
}

func (e *engine) BindRoutes() error {
	for _, route := range e.routers {
		if route == nil {
			continue
		}
		return e.router.BindRoute(route, e.middlewares)
	}
	return nil
}

func (e *engine) AddRoute(router *router.Route) {
	if router != nil { // 在添加路由前进行空指针检查
		e.routers = append(e.routers, router)
	}
}

func (e engine) Run(addr string) {
	err := fasthttp.ListenAndServe(addr, e.router.Handle)
	if err != nil {
		panic(err)
	}
}
