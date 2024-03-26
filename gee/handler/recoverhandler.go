package handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

// ErrorMiddleware 是一个错误处理中间件，它会捕获并处理传入的请求处理器中的错误。
func ErrorMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				// 处理panic情况
				errorMsg := fmt.Sprintf("Internal Server Error: %v", err)
				handleError(ctx, errorMsg, fasthttp.StatusInternalServerError)
			}
		}()

		next(ctx)

		// 检查ctx是否有设置错误
		if ctx.Response.StatusCode() >= fasthttp.StatusBadRequest {
			// 已经设置了错误状态码，不需要额外处理
			return
		}

		// 如果有错误信息，则处理错误
		if errStr := ctx.UserValue("error"); errStr != nil {
			errorMsg := errStr.(string)
			statusCode := fasthttp.StatusInternalServerError // 默认状态码
			// 根据实际情况，可以从ctx获取更精确的状态码
			handleError(ctx, errorMsg, statusCode)
		}
	}
}

// handleError 用于发送带有错误信息和状态码的HTTP响应
func handleError(ctx *fasthttp.RequestCtx, errorMsg string, statusCode int) {
	ctx.SetStatusCode(statusCode)
	ctx.SetBodyString(errorMsg)
}
