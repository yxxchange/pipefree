package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/operator"
)

const (
	routeGroup = "/operator"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.GET("/namespace/:namespace/name/:name/kind/:kind", Watch)
	}
}

func Watch(c *gin.Context) {
	var req WatchReq
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseError(c, -1, err.Error())
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
			common.ResponseOk(c, event)
		case <-c.Request.Context().Done():
			log.Infof("watch context done for prefix %s", req.KeyPrefix())
			eventCh.Close()                             // 关闭事件通道
			common.ResponseOk(c, operator.CloseEvent()) // 返回关闭事件
			return                                      // 处理上下文取消
		}
	}
}
