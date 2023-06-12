package jenkins

import (
	"devopscenter/model"
	"devopscenter/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

//func Build(c *gin.Context) {
//	response := model.Res{
//		Code:    201,
//		Message: "successful",
//		Data:    nil,
//	}
//	name := c.Query("job_name")
//	result := service.CheckJob(name)
//	if result == false {
//		response.Message = "Jenkins Job Not Exist"
//		response.Data = result
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	data := model.JenkinsTemplate{}
//	err := c.ShouldBindJSON(&data)
//	if err != nil {
//		response.Message = "Json Paras Failed"
//		response.Data = err
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	_, err1 := service.BuildJob(name, &data)
//	if err1 != nil {
//		response.Message = "Build Failed"
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	result1, err2 := service.IdJob(name)
//	if err2 != nil {
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	response.Data = result1
//	c.JSON(http.StatusCreated, response)
//}

func BuildV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	data := model.JenkinsTemplate{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Message = "Json Paras Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	//name := data.Project + "_" + data.Env
	name := data.Language + "_Template"
	// 检查Jenkins job 是否存在
	if result := service.CheckJob(name); result == false {
		response.Message = "Jenkins Job " + name + " Not Exist"
		response.Data = result
		c.JSON(http.StatusOK, response)
		return
	}
	// 构建指定job
	if _, err := service.BuildJob(name, &data); err != nil {
		response.Message = "Build Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	// 获取构建ID
	result, err2 := service.IdJob(name)
	if err2 != nil {
		c.JSON(http.StatusOK, response)
		return
	}
	// 序列化参数记录构建信息，存入数据库
	if marshalData, err := json.Marshal(data); err != nil {
		response.Message = "Json Marshal Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	} else {
		if !service.RecordBuildInfo(&data, string(marshalData), int(result)) {
			response.Data = false
			response.Message = "RecordBuildInfo Failed"
			c.JSON(http.StatusOK, response)
			return
		}
	}
	response.Data = result + 1
	c.JSON(http.StatusOK, response)
}

func GetJobId(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	name := c.Query("job-name")
	id, err := service.IdJob(name)
	response.Data = id
	if err != nil {
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func Status(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	result := service.StatusJob(c.Query("job-name"))
	response.Data = result
	if result == "" {
		response.Message = "Build Ing"
	} else if result == "ABORTED" {
		response.Message = "Build Aborted"
	} else if result == "SUCCESS" {
		response.Message = "Build Successful"
	} else if result == "FAILURE" {
		response.Message = "Build Failed"
	}
	c.JSON(http.StatusOK, response)
}
