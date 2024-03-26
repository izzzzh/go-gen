package main

import (
	"go-gen/gee"
	"go-gen/gee/router"
	"go-gen/internal/handler"
	"net/http"
)

func main() {
	engine := gee.NewEngine()
	engine.AddRoute(&router.Route{
		Method:  http.MethodGet,
		Path:    "/hello/:id",
		Handler: handler.HelloHandler(),
	})
	engine.Start()
}
