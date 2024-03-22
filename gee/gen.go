package gee

import (
	"github.com/valyala/fasthttp"
	"go-gen/gee/router"
)

func Start(r *router.Router) {
	err := r.BindRoute()
	if err != nil {
		panic(err)
	}
	Run(":8080", r)
}

func Run(addr string, r *router.Router) {
	err := fasthttp.ListenAndServe(addr, r.Handle)
	if err != nil {
		panic(err)
	}
}
