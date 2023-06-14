package router

import (
	"devopscenter/controller/jenkins"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
)

func JenkinsRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.POST("/jenkins/job", middleware.JwtAuth(), jenkins.BuildV2)
		//api.OPTIONS("/jenkins/job", func(c *gin.Context) {
		//	c.Status(http.StatusOK)
		//})
		api.GET("/jenkins/job-id", middleware.JwtAuth(), jenkins.GetJobId)
		//api.GET("/jenkins/job-status", jenkins.Status)
		api.GET("/jenkins/:name/:id", middleware.JwtAuth(), jenkins.StatusV2)
	}
}
