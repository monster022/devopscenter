package router

import (
	"devopscenter/controller/harbor"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
)

func HarborRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.GET("/harbor", middleware.JwtAuth(), harbor.List)
	}
}
