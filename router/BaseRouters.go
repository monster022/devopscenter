package router

import (
	"devopscenter/controller/authority"
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
	c.GET("/base", base.Get)
	c.POST("/base/json", base.PostJson)
	c.POST("/base/from-data", base.PostForm)
	c.DELETE("/base", base.Delete)
	c.PATCH("/base", base.Patch)

	api := c.Group("/devops/")
	{
		api.GET("/menu/:name", authority.List)

		api.POST("/login", login.Auth)
		api.POST("/password", middleware.JwtAuth(), login.ModifyPassword)
	}
}
