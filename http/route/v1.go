package route

import (
	"github.com/gin-gonic/gin"
	"github.com/yxxchange/pipefree/http/api/pipe_cfg"
	"github.com/yxxchange/pipefree/http/api/pipe_exec"
	"github.com/yxxchange/pipefree/http/api/pipe_perm"
	"github.com/yxxchange/pipefree/http/api/pipe_watch"
)

const (
	v1 = "/api/v1"
)

func RegisterV1Routes(router *gin.Engine) {
	group := router.Group(v1)
	pipe_cfg.RegisterV1(group)
	pipe_exec.RegisterV1(group)
	pipe_perm.RegisterV1(group)
	pipe_watch.RegisterV1(group)
}
