package handler

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

func LoggingMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		startTime := time.Now()

		// 请求前打印日志
		requestLog := fmt.Sprintf("[Start] Request at %s, Method: %s, Path: %s\n",
			startTime.Format(time.RFC3339),
			ctx.Method(),
			ctx.Path(),
		)
		fmt.Print(requestLog)

		// 执行下一个中间件或处理函数
		next(ctx)

		// 请求完成后打印日志
		endTime := time.Now()
		elapsed := endTime.Sub(startTime)
		responseLog := fmt.Sprintf("[End] Response after %s, Status Code: %d, Path: %s\n",
			elapsed,
			ctx.Response.StatusCode(),
			ctx.Path(),
		)
		fmt.Print(responseLog)
	}
}
