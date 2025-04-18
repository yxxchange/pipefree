package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/bridge/pool"
	"github.com/yxxchange/pipefree/pkg/http/utils"
	"sync"
)

type WatchParam struct {
	Namespace string `uri:"namespace" binding:"required"`
	Name      string `uri:"name" binding:"required"`
	Kind      string `query:"kind" binding:"required"`
}

type WatchCtx struct {
	mu    sync.Mutex
	ctx   *gin.Context
	param WatchParam
}

func Integrate(ctx *gin.Context, param WatchParam) pool.Connection {
	return &WatchCtx{
		param: param,
		ctx:   setChunkedConnection(ctx),
	}
}

var _ pool.Connection = &WatchCtx{}

func (ctx *WatchCtx) Write(data interface{}) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	utils.ResponseOKWithInfo(ctx.ctx, data)
	return nil
}

func (ctx *WatchCtx) ConnectionType() string {
	return ctx.param.Kind
}

func setChunkedConnection(c *gin.Context) *gin.Context {
	c.Request.Header.Set("Connection", "keep-alive")
	c.Request.Header.Set("Transfer-Encoding", "chunked")
	c.Request.Header.Set("Content-Type", "application/json")
	return c
}
