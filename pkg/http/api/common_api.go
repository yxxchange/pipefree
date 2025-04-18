package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/view"
)

func HealthCheck(ctx *gin.Context) {
	// TODO: do more
	view.ResponseOK(ctx)
}
