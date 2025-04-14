package http

import (
	"github.com/gin-gonic/gin"
	"time"
)

const (
	MetricTimeKey = "http::api::start"
)

func MetricTimeCost(ctx *gin.Context) {
	ctx.Set(MetricTimeKey, time.Now())
	ctx.Next()
}
