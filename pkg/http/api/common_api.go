package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/http/utils"
)

func HealthCheck(ctx *gin.Context) {
	// TODO: do more
	utils.ResponseOK(ctx)
}
