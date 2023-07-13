package router

import (
	"devopscenter/controller/machine"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
)

func MachineRegister(c *gin.Engine) {
	api := c.Group("/devops/")
	{
		api.GET("/machine", middleware.JwtAuth(), machine.ListV2)
		api.GET("/machine/password", middleware.JwtAuth(), machine.Password)
		api.POST("/machine", middleware.JwtAuth(), machine.Create)
		api.DELETE("/machine", middleware.JwtAuth(), machine.Remove)
		api.PUT("/machine", middleware.JwtAuth(), machine.Update)
		api.PATCH("/machine/:id", middleware.JwtAuth(), machine.PatchName)
	}
}
