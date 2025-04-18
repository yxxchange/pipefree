package http

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yxxchange/pipefree/pkg/http/api"
)

func Launch() {
	r := gin.Default(initMiddleware, initRoute)
	err := r.Run(viper.GetString("http.port"))
	if err != nil {
		panic(err)
	}
}

func initMiddleware(r *gin.Engine) {
	r.Use(MetricTimeCost)
}

func initRoute(r *gin.Engine) {
	r.GET("/health", api.HealthCheck)
	v1(r.Group("/api/v1"))
	v2(r.Group("/api/v2"))
	v3(r.Group("/api/v3"))
}

func v1(r *gin.RouterGroup) {
	pipe := r.Group("/pipe/namespace/:namespace/name/:name")
	{
		pipe.GET("/watch", api.Watch)
		pipe.GET("/list", api.List)
	}

}

func v2(r *gin.RouterGroup) {
	// TODO: implement v2 if needed
}

func v3(r *gin.RouterGroup) {
	// TODO: implement v3 if needed
}
