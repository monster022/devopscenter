package router

import (
	"devopscenter/controller/base"
	"devopscenter/controller/login"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PostData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func BaseRegister(c *gin.Engine) {
	c.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
	api := c.Group("/devops/", middleware.JwtAuth())
	{
		api.GET("/base", base.Get)
		api.POST("/base/json", base.PostJson)
		api.POST("/base/from-data", base.PostForm)
		api.DELETE("/base", base.Delete)
		api.PATCH("/base", base.Patch)

		api.POST("/login", login.Auth)
	}
}
