package router

import (
	"devopscenter/controller/jenkins"
	"github.com/gin-gonic/gin"
)

func JenkinsRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.POST("/jenkins/job", jenkins.BuildV2)
		//api.OPTIONS("/jenkins/job", func(c *gin.Context) {
		//	c.Status(http.StatusOK)
		//})
		api.GET("/jenkins/job-id", jenkins.GetJobId)
		api.GET("/jenkins/job-status", jenkins.Status)
	}
}
