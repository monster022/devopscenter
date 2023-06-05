package gitlab

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	response := model.Res{
		Code:    201,
		Message: "successful",
		Data:    nil,
	}
	data := model.Application{}
	err := c.ShouldBindJSON(&data)
	if err != nil {
		response.Message = "Json Parse Failed"
		c.JSON(http.StatusCreated, response)
		return
	}
	//
	result, pid, repo, err1 := service.SearchName(data.Name)
	if result == false || err1 != nil {
		response.Message = "Project Not Exist"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	project := model.Project{}
	project.ProjectId = pid
	project.ProjectName = data.Name
	project.ProjectRepo = repo
	project.ProjectStatus = 1
	result1 := project.Insert()
	response.Data = result1
	if result1 == false {
		response.Message = "Project Exist"
		c.JSON(http.StatusCreated, response)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func Patch(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	projectStatus := c.Query("project_status")
	var Id = c.Query("id")
	status, err1 := strconv.Atoi(projectStatus)
	id, err2 := strconv.Atoi(Id)
	if err1 != nil || err2 != nil {
		response.Message = "Type Convert Failed"
		response.Data = err1
		c.JSON(http.StatusOK, response)
		return
	}
	project := model.Project{}

	result := project.Patch(id, status)
	response.Data = result
	if result == false {
		response.Message = "Project_Status Charge Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func List(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	projectPage := c.Query("page")
	projectSize := c.Query("size")
	page, err1 := strconv.Atoi(projectPage)
	size, err2 := strconv.Atoi(projectSize)
	if err1 != nil || err2 != nil {
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	project := model.Project{}
	data := project.List(page, size)
	total := project.Count()
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    data,
		"total":   total,
	})
}

func BranchList(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	projectId := c.Query("id")
	id, err1 := strconv.Atoi(projectId)
	if err1 != nil {
		response.Message = "Type Convert Failed"
		response.Data = err1
		c.JSON(http.StatusOK, response)
		return
	}
	data, err2 := service.BranchList(id)
	response.Data = data
	if err2 != nil {
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func Search(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	project := c.Query("project")
	if project == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	result, pid, repo, err := service.SearchName(project)
	if result == false || err != nil {
		response.Message = "Project Not Exist"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	data := make(map[string]interface{})
	data["id"] = pid
	data["repo"] = repo
	data["name"] = project
	response.Message = "Project Exist"
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func SearchAll(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	project := c.Query("project")
	if project == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	data, err := service.SearchAll(project)
	response.Data = data
	if err != nil {
		response.Message = "Project Not Exist"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}
