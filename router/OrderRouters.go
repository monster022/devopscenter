package router

import (
	"devopscenter/controller/order"
	"github.com/gin-gonic/gin"
)

func OrderRegister(c *gin.Engine) {
	api := c.Group("/devops/")
	{
		api.GET("/order", order.List)
		api.GET("/order/tackle/:name", order.ListTackleName)
		api.GET("/order/tackle/:name/total", order.ListTackleNameTotal)
		api.GET("/order/submit/:name", order.ListSubmitName)

		api.PATCH("/order/:id/:status", order.PatchOrder)
		api.POST("/order", order.Create)
		api.POST("/order/:id/:status", order.PostOrder)
	}
}
