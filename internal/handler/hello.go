package handler

import (
	"github.com/valyala/fasthttp"
	"go-gen/gee/httpx"
	"go-gen/internal/logic"
	"go-gen/internal/types"
	"net/http"
)

func HelloHandler() fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		var req types.HelloReq
		if err := httpx.Parse(ctx, &req); err != nil {
			ctx.Error(err.Error(), http.StatusBadRequest)
			return
		}

		l := logic.NewHelloLogic(ctx)
		ret, err := l.Hello(&req)
		if err != nil {
			ctx.Error(err.Error(), http.StatusBadRequest)
		} else {
			httpx.JSON(ctx, http.StatusOK, ret)
		}
	}
}
