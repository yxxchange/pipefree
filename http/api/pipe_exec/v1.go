package pipe_exec

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/pipe_exec"
)

const (
	routeGroup = "/pipe_exec"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.POST("/:pipe_id", Run)
	}
}

func Run(c *gin.Context) {
	var req PipeExecReqParam
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseError(c, -1, "invalid request parameters")
		return
	}
	if req.PipeId <= 0 {
		common.ResponseError(c, -1, "pipe id is required")
		return
	}
	err := pipe_exec.NewService(c).Run(req.PipeId)
	if err != nil {
		common.ResponseError(c, pipe_exec.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, "pipe execution started successfully")
}
