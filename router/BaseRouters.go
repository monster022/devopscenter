package router

import (
	"devopscenter/controller/base"
	"devopscenter/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PostData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func BaseRegister(c *gin.Engine) {
	c.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
	api := c.Group("/devops/")
	{
		api.POST("/api/postData", func(c *gin.Context) {
			fmt.Println("我进来了 setp 1")
			var postData PostData
			if err := c.BindJSON(&postData); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fmt.Println(postData)

			c.JSON(200, gin.H{
				"message": "ok",
			})
		})
		api.GET("/base", base.Get)
		api.POST("/base/json", base.PostJson)
		api.POST("/base/from-data", base.PostForm)
		api.DELETE("/base", base.Delete)
		api.PATCH("/base", base.Patch)
	}
}
