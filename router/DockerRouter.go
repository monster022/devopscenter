package router

import (
	"devopscenter/controller/docker"
	"github.com/gin-gonic/gin"
)

func DockerRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.POST("/container", docker.Create)
		api.GET("/container/machine/:id", docker.Machine)
	}
}
