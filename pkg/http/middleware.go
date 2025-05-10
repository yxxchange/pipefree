package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/http/utils"
	"time"
)

func MetricTimeCost(ctx *gin.Context) {
	ctx.Set(utils.MetricTimeKey, time.Now())
	ctx.Next()
}
