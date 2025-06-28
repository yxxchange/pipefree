package pipe_cfg

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/helper/log"
	"github.com/yxxchange/pipefree/http/common"
	"github.com/yxxchange/pipefree/service/pipe_cfg"
)

const (
	routeGroup = "/pipe_cfg"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.GET("/:pipeId", Get)
		group.POST("", Create)
	}
}

func Get(c *gin.Context) {
	var req PipeReqParam
	err := c.ShouldBindUri(&req)
	if err != nil {
		log.Errorf("Get pipe configuration failed, invalid request parameters: %v", err)
		common.ResponseError(c, -1, "Invalid request parameters")
		return
	}
	cfg, err := pipe_cfg.NewService(c).GetById(req.PipeId)
	if err != nil {
		common.ResponseError(c, pipe_cfg.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, Convert(*cfg))
}

func Create(c *gin.Context) {
	var req PipeReqParam
	if err := c.ShouldBind(&req); err != nil {
		log.Errorf("Create pipe configuration failed, invalid request parameters: %v", err)
		common.ResponseError(c, -1, "Invalid request parameters")
		return
	}
	if err := pipe_cfg.NewService(c).Create(req.View.PipeCfg, req.View.NodeCfgList); err != nil {
		common.ResponseError(c, pipe_cfg.ErrorCode, err.Error())
		return
	}
	common.ResponseOk(c, "Pipe configuration created successfully")
}
