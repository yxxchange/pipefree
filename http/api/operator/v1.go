package operator

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/http/common"
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

	common.ResponseOk(c, nil)
}
