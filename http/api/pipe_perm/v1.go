package pipe_perm

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/pipe_perm"
)

const (
	routeGroup = "/pipe_perm"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.POST("/space/:space", CreatePipeSpace)
		group.POST("/namespace/:namespace", CreateNodeNamespace)
		group.POST("/apply", ApplyPipePermission)
	}
}

func CreatePipeSpace(c *gin.Context) {
	var req PipePermReqParam
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseError(c, -1, "invalid request parameters")
		return
	}
	if req.Space == "" {
		common.ResponseError(c, -1, "pipe space is required")
		return
	}
	err := pipe_perm.NewService(c).CreatePipeSpace(req.Space)
	if err != nil {
		common.ResponseError(c, pipe_perm.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, "pipe space created successfully")
}

func CreateNodeNamespace(c *gin.Context) {
	var req PipePermReqParam
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseError(c, -1, "invalid request parameters")
		return
	}
	if req.Namespace == "" {
		common.ResponseError(c, -1, "node namespace is required")
		return
	}
	err := pipe_perm.NewService(c).CreateNodeNamespace(req.Namespace)
	if err != nil {
		common.ResponseError(c, pipe_perm.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, "node namespace created successfully")
}

func ApplyPipePermission(c *gin.Context) {
	var req PipePermReqParam
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, -1, "invalid request parameters")
		return
	}
	if req.Space == "" || req.Namespace == "" {
		common.ResponseError(c, -1, "pipe space and node namespace are required")
		return
	}
	err := pipe_perm.NewService(c).CreatePermissionItem(req.Space, req.Namespace)
	if err != nil {
		common.ResponseError(c, pipe_perm.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, "pipe permission subscribed successfully")
}
