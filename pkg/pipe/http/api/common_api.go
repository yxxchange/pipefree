package api

import (
	"github.com/gin-gonic/gin"
	"pipefree/pkg/pipe/http"
)

func HealthCheck(ctx *gin.Context) {
	http.ResponseOK(ctx)
}
