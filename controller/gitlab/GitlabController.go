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
		Code:    20000,
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
	// gitlab中查找该项目
	result, pid, repo, err1 := service.SearchName(data.Name)
	if result == false || err1 != nil {
		response.Message = "Project Not Exist"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	// 数据库中添加项目
	project := model.Project{}
	project.ProjectId = pid
	project.ProjectName = data.Name
	project.ProjectRepo = repo
	project.ProjectStatus = 1
	project.Language = data.Language
	project.BuildPath = data.BuildPath
	project.PackageName = data.PackageName
	project.AliasName = data.AliasName
	result1 := project.Insert()
	response.Data = result1
	if result1 == false {
		response.Message = "Project Exist"
		c.JSON(http.StatusCreated, response)
		return
	}
	c.JSON(http.StatusCreated, response)
}

func StatusPatch(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	status, err1 := strconv.Atoi(c.Query("status"))
	id, err2 := strconv.Atoi(c.Query("id"))
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
		response.Message = "Project_Status Change Failed"
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

func Delete(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Message = "id type convert failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	data := model.Project{}
	result := data.Delete(id)
	response.Data = result
	if result == false {
		response.Message = "filed delete failed"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func EditPatch(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	type body struct {
		BuildPath   string `json:"build_path"`
		PackageName string `json:"package_name"`
	}
	json := body{}
	if err := c.ShouldBindJSON(&json); err != nil {
		response.Message = "Json Parse Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	project := model.Project{}
	if result := project.Edit(c.Param("name"), json.BuildPath, json.PackageName); result == false {
		response.Message = "Modify Failed"
		response.Data = result
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
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
		response.Data = err1.Error()
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

func CommitMessage(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	commit, err := service.CommitByIdAndBranch(c.Param("pid"), c.Query("branch"))
	if err != nil {
		response.Message = "Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	message := struct {
		AuthorName string `json:"authorName"`
		Message    string `json:"message"`
	}{
		AuthorName: commit.AuthorName,
		Message:    commit.Title,
	}
	response.Data = message
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
	if err != nil {
		response.Message = "Internal Server Error"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	if !result {
		response.Message = "Project Not Exist"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	data := map[string]interface{}{
		"id":   pid,
		"repo": repo,
		"name": project,
	}
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

func ListDetail(c *gin.Context) {
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
	project := model.ProjectDetail{}
	result := project.List(c.Param("name"), page, size)
	response.Data = result
	c.JSON(http.StatusOK, response)
}

func ListDeployDetail(c *gin.Context) {
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
	deployProject := model.DeployProjectDetail{}
	result, err := deployProject.List(c.Param("name"), page, size)
	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = result
	c.JSON(http.StatusOK, response)
}
