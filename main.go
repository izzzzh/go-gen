package main

import (
	"go-gen/gee"
	"go-gen/gee/router"
	"go-gen/internal/handler"
	"net/http"
)

func main() {
	r := router.New()
	r.AddRoute(&router.Route{
		Method:  http.MethodGet,
		Path:    "/hello/:id",
		Handler: handler.HelloHandler(),
	})

	gee.Start(r)
}
