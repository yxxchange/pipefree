package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/pkg/http"
)

func HealthCheck(ctx *gin.Context) {
	http.ResponseOK(ctx)
}
