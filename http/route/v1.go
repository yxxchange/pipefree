package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/http/api/pipe_cfg"
)

const (
	v1 = "/api/v1"
)

func RegisterV1Routes(router *gin.Engine) {
	group := router.Group(v1)
	pipe_cfg.RegisterV1(group)
}
