package router

import (
	"devopscenter/controller/gitlab"
	"devopscenter/middleware"
	"github.com/gin-gonic/gin"
)

func GitlabRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.POST("/project", middleware.JwtAuth(), gitlab.Create)
		api.PATCH("/project", middleware.JwtAuth(), gitlab.StatusPatch)
		api.GET("/project", middleware.JwtAuth(), gitlab.List)
		api.DELETE("/project/:id", middleware.JwtAuth(), gitlab.Delete)
		api.PATCH("/project/:name", middleware.JwtAuth(), gitlab.EditPatch)

		api.GET("/project/branch", middleware.JwtAuth(), gitlab.BranchList)
		api.GET("/project/:pid/commit", gitlab.CommitMessage)

		api.GET("/project/search", middleware.JwtAuth(), gitlab.Search)
		api.GET("/project/search/all", middleware.JwtAuth(), gitlab.SearchAll)

		api.GET("/project/detail/:name", middleware.JwtAuth(), gitlab.ListDetail)

	}
}
