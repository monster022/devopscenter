package jenkins

import (
	"devopscenter/configuration"
	"devopscenter/model"
	"devopscenter/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
	name := data.Language + "_Template"
	// 检查Jenkins job 是否存在
	if err := service.CheckJob(name); err != nil {
		response.Message = "Jenkins Job " + name + " Not Exist"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}

	// 构建指定job
	if _, err := service.BuildJob(name, &data); err != nil {
		response.Message = "Build Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}

	// 获取构建ID
	result, err := service.IdJob(name)
	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}

	// 序列化参数记录构建信息，存入数据库
	if marshalData, err := json.Marshal(data); err != nil {
		response.Message = "Json Marshal Failed"
		response.Data = err.Error()
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

func StatusV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	url := configuration.Configs.JenkinsUrl + "job/" + c.Param("name") +
		"/" + c.Param("id") + "/api/json"
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(configuration.Configs.JenkinsUsername, configuration.Configs.JenkinsPassword)
	if err != nil {
		response.Data = err
		response.Message = "Http Request Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	resp, err1 := (&http.Client{}).Do(req)
	if err1 != nil {
		response.Data = err1
		response.Message = "HttpClient Do Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	var data struct {
		Building bool   `json:"building"`
		Result   string `json:"result"`
	}
	if ok := json.NewDecoder(resp.Body).Decode(&data); ok != nil {
		response.Message = "NewDecoder Failed"
		response.Data = ok
		c.JSON(http.StatusOK, response)
		return
	}
	project := model.ProjectDetail{}
	if data.Result == "" {
		data.Result = "ING"
	} else {
		jobId, err2 := strconv.Atoi(c.Param("id"))
		if err2 != nil {
			response.Data = err2
			c.JSON(http.StatusOK, response)
			return
		}
		if err3 := project.Update(c.Param("name"), data.Result, jobId); err3 != true {
			response.Message = "Database Update Failed"
			c.JSON(http.StatusOK, response)
			return
		}
	}
	response.Data = data.Result
	c.JSON(http.StatusOK, response)
}
