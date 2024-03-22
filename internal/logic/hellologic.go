package logic

import (
	"github.com/valyala/fasthttp"
	"go-gen/internal/types"
)

type HelloLogic struct {
	ctx *fasthttp.RequestCtx
}

func NewHelloLogic(ctx *fasthttp.RequestCtx) *HelloLogic {
	return &HelloLogic{
		ctx: ctx,
	}
}

func (h *HelloLogic) Hello(req *types.HelloReq) (*types.HelloResp, error) {
	ret := &types.HelloResp{
		Code:    200,
		Message: "ok",
		Data:    "hello world",
	}
	return ret, nil
}
