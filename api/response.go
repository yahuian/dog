package api

import (
	"dog/errcode"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func handleError(ctx *gin.Context, err error) {
	code, msg := errcode.Parse(err)

	// 如果某个错误没有加 errcode 统一认为服务端内部错误，并隐藏错误信息
	if code == 0 && msg == "" {
		code = http.StatusInternalServerError
	}

	// 服务端内部错误不返回详细信息
	if code == http.StatusInternalServerError {
		resp := Resp{
			Code: http.StatusInternalServerError,
			Msg:  "服务器开小差了",
		}
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 其他自定义错误，不返回错误链上的其他信息，仅返回用户友好的提示信息
	resp := Resp{
		Code: code,
		Msg:  msg,
	}
	ctx.JSON(errcode.Code2Status(code), resp)
}

func Success(ctx *gin.Context, msg string, data any) {
	resp := Resp{
		Code: http.StatusOK,
		Msg:  msg,
		Data: data,
	}

	ctx.JSON(http.StatusOK, resp)
}
