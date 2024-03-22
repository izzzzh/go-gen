package httpx

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
)

func JSON(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(code)

	if ret, err := json.Marshal(obj); err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
	} else {
		ctx.Write(ret)
	}
}
