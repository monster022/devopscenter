package router

import (
	"devopscenter/controller/machine"
	"github.com/gin-gonic/gin"
)

func MachineRegister(c *gin.Engine) {
	api := c.Group("/devops/")
	{
		api.GET("/machine", machine.List)
		api.GET("/machine/password", machine.Password)
		api.POST("/machine", machine.Create)
		api.DELETE("/machine", machine.Remove)
		api.PUT("/machine", machine.Update)
	}
}
