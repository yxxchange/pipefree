package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Code string

const (
	MetricTimeKey = "http::api::start"
)

const (
	CodeSuccess Code = "0"
	CodeError   Code = "-1"
)

type Response struct {
	Code   Code        `json:"code"`
	ErrMsg string      `json:"err_msg"`
	Info   interface{} `json:"info"`
	Cost   string      `json:"cost"` // time cost(ms)
}

func ResponseOKWithInfo(ctx *gin.Context, info interface{}) {
	resp := Response{
		Code:   CodeSuccess,
		ErrMsg: "",
		Info:   info,
		Cost:   calculateCost(ctx),
	}
	ctx.JSON(200, resp)
}

func ResponseOK(ctx *gin.Context) {
	resp := Response{
		Code:   CodeSuccess,
		ErrMsg: "",
		Info:   nil,
		Cost:   calculateCost(ctx),
	}
	ctx.JSON(200, resp)
}

func ResponseError(ctx *gin.Context, err error) {
	resp := Response{
		Code:   CodeError,
		ErrMsg: err.Error(),
		Info:   nil,
		Cost:   calculateCost(ctx),
	}
	ctx.JSON(200, resp)
}

func ResponseErrorWithCode(ctx *gin.Context, code Code, err error) {
	resp := Response{
		Code:   code,
		ErrMsg: err.Error(),
		Info:   nil,
		Cost:   calculateCost(ctx),
	}
	ctx.JSON(200, resp)
}

func calculateCost(ctx *gin.Context) string {
	v, ok := ctx.Get(MetricTimeKey)
	if !ok {
		return ""
	}
	start, ok := v.(time.Time)
	if !ok {
		return ""
	}
	cost := time.Since(start).Milliseconds()
	return fmt.Sprintf("%dms", cost)
}

type BindFunc func(interface{}) error

func BindFlow(obj interface{}, binds ...BindFunc) error {
	for _, bind := range binds {
		if err := bind(obj); err != nil {
			return err
		}
	}
	return nil
}
