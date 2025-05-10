package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/http/internal"
	"github.com/yxxchange/pipefree/pkg/http/utils"
	"github.com/yxxchange/pipefree/pkg/pipe/model"
)

func CreatePipe(c *gin.Context) {
	var n model.Node
	err := utils.BindFlow(&n, c.ShouldBindYAML)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	err = internal.CreatePipe(context.TODO(), n.ToPipeCfg())
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseOK(c)
}

func RunPipe(c *gin.Context) {
	type RunPipeReq struct {
		Id string `form:"id" bind:"required"`
	}
	var req RunPipeReq
	err := utils.BindFlow(&req, c.ShouldBindQuery)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	err = internal.RunPipe(context.TODO(), req.Id)
	if err != nil {
		utils.ResponseError(c, err)
		return
	}
	utils.ResponseOK(c)
}
