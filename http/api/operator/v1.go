package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/operator"
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
	eventCh := operator.NewService(c).Watch(req.KeyPrefix())
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
				// TODO: 关闭etcd连接
				return // 处理写入错误
			}
			flusher.Flush()
		case <-c.Request.Context().Done():
			log.Infof("watch context done for prefix %s", req.KeyPrefix())
			eventCh.Close()                             // 关闭事件通道
			common.ResponseOk(c, operator.CloseEvent()) // 返回关闭事件
			return                                      // 处理上下文取消
		}
	}
}
