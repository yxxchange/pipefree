package pipe_watch

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/pipe_watch"
	"io"
	"net/http"
)

const (
	routeGroup = "/operator"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.GET("/namespace/:namespace/kind/:kind", Watch)
	}
}

func Watch(c *gin.Context) {
	var req WatchReq
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseError(c, -1, err.Error())
		return
	}
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		common.ResponseError(c, -1, "response writer does not support flushing")
		return
	}
	operator := pipe_watch.NewOperatorCtx(req.KeyPrefix(), uuid.New().String())
	watchService := pipe_watch.NewService(c)
	eventCh := watchService.Watch(operator)
	defer func() {
		watchService.UnWatch(operator)
	}()
	for {
		select {
		case event, ok := <-eventCh.Ch():
			if !ok {
				log.Warnf("event channel closed for prefix %s", req.KeyPrefix())
				return // 处理通道关闭
			}
			_, err := io.Writer(c.Writer).Write(event)
			if err != nil {
				log.Errorf("failed to write event for prefix %s: %v", req.KeyPrefix(), err)
				return // 处理写入错误
			}
			flusher.Flush()
		case err := <-eventCh.ErrCh(): // 连接断开
			if err != nil {
				log.Errorf(err.Error())
				return
			}
		case <-c.Request.Context().Done():
			log.Info("http context done, exit!")
			common.ResponseOk(c, pipe_watch.CloseEvent()) // 返回关闭事件
			return                                        // 处理上下文取消
		}
	}
}
