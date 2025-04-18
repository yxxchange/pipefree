package pool

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/view"
	"sync"
)

const (
	ChunkedEmptyBlock = "0\r\n\r\n"
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

var pool *ConnectionPool

type ConnectionPool struct {
	mutex sync.Mutex
	pool  map[string][]*WatchCtx
}

func GetPool() *ConnectionPool {
	if pool == nil {
		pool = &ConnectionPool{
			pool: make(map[string][]*WatchCtx),
		}
	}
	return pool
}

func (c *ConnectionPool) Store(ctx *WatchCtx) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pool[ctx.param.Kind] = append(c.pool[ctx.param.Kind], ctx)
}

func (c *ConnectionPool) Consume(event view.Event) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	connArr := c.pool[event.Kind]
	wg := sync.WaitGroup{}
	wg.Add(len(connArr))
	for _, conn := range connArr {
		go func(conn *WatchCtx) {
			defer wg.Done()
			view.ResponseOKWithInfo(conn.ctx, event)
		}(conn)
	}
	wg.Wait()
}

func (c *ConnectionPool) Transport(event view.Event) {

}

func Register(ctx *gin.Context, param WatchParam) {
	setChunkedConnection(ctx)
	watchCtx := &WatchCtx{
		ctx:   ctx,
		param: param,
	}
	if pool == nil {
		pool = &ConnectionPool{
			pool: make(map[string][]*WatchCtx),
		}
	}
	pool.Store(watchCtx)
	return
}

func setChunkedConnection(c *gin.Context) {
	c.Request.Header.Set("Connection", "keep-alive")
	c.Request.Header.Set("Transfer-Encoding", "chunked")
	c.Request.Header.Set("Content-Type", "application/json")
}
