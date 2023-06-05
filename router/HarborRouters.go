package router

import (
	"devopscenter/controller/harbor"
	"github.com/gin-gonic/gin"
)

func HarborRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.GET("/harbor", harbor.List)
	}
}
