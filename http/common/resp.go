package common

import "github.com/gin-gonic/gin"

type Response struct {
	Code   int         `json:"code"`    // 状态码, 0表示成功, 非0表示失败
	ErrMsg string      `json:"err_msg"` // 错误信息, 成功时为空
	Info   interface{} `json:"info"`    // 附加信息, 成功时包含数据
}

func ResponseOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Response{
		Code:   0,
		ErrMsg: "",
		Info:   data,
	})
}

func ResponseError(ctx *gin.Context, code int, errMsg string) {
	ctx.JSON(200, Response{
		Code:   code,
		ErrMsg: errMsg,
		Info:   nil,
	})
}
