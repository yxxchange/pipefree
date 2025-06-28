package operator

import "github.com/gin-gonic/gin"

const (
	routeGroup = "/operator"
)

func RegisterV1(router *gin.RouterGroup) {
	group := router.Group(routeGroup)
	{
		group.GET("namespace/:namespace/name/:name", ListAndWatch)
	}
}

func ListAndWatch(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	if namespace == "" || name == "" {
		c.JSON(400, gin.H{"error": "namespace and name are required"})
		return
	}

	// Here you would typically call a service to list and watch the operators.
	// For now, we'll just return a mock response.
	response := gin.H{
		"namespace": namespace,
		"name":      name,
		"status":    "active",
	}

	c.JSON(200, response)
}
