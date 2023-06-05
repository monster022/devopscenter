package router

import (
	"devopscenter/controller/gitlab"
	"github.com/gin-gonic/gin"
)

func GitlabRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.POST("/project", gitlab.Create)
		api.PATCH("/project", gitlab.Patch)
		api.GET("/project", gitlab.List)

		api.GET("/project/branch", gitlab.BranchList)

		api.GET("/project/search", gitlab.Search)
		api.GET("/project/search/all", gitlab.SearchAll)
	}
}
