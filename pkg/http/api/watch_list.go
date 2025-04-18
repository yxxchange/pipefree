package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/bridge/pool"
	"github.com/yxxchange/pipefree/pkg/http/utils"
	"github.com/yxxchange/pipefree/pkg/view"
	"strings"
)

func Watch(ctx *gin.Context) {
	if !isKeepAlive(ctx) {
		view.ResponseError(ctx, fmt.Errorf("must be a chunked keep-alive connection"))
		return
	}
	var watchParam pool.WatchParam
	err := utils.BindFlow(&watchParam, ctx.ShouldBindUri, ctx.ShouldBindQuery)
	if err != nil {
		view.ResponseError(ctx, err)
		return
	}
	pool.Register(ctx, watchParam)
	view.ResponseOK(ctx)
}

func isKeepAlive(c *gin.Context) bool {
	// 默认 HTTP/1.1 启用 Keep-Alive，除非显式设置 Connection: close
	if c.Request.ProtoAtLeast(1, 1) {
		connectionHeader := c.Request.Header.Get("Connection")
		if strings.EqualFold(connectionHeader, "close") {
			return false
		}
		return true
	}
	// HTTP/1.0 需要显式设置 Connection: keep-alive
	connectionHeader := c.Request.Header.Get("Connection")
	return strings.EqualFold(connectionHeader, "keep-alive")
}

func List(ctx *gin.Context) {
	// TODO: implement list
}
