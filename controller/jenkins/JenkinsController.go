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

//func Status(c *gin.Context) {
//	response := model.Res{
//		Code:    20000,
//		Message: "successful",
//		Data:    nil,
//	}
//
//	url := "http://jenkins.chengdd.cn/job/" + c.Query("name") + "/" + c.Query("id") + "/api/json"
//	fmt.Println(url)
//	req, err := http.NewRequest("GET", url, nil)
//	req.SetBasicAuth("yen", "1qaz@WSX")
//	if err != nil {
//		response.Message = "Request Failed"
//		response.Data = err
//		c.JSON(http.StatusOK, response)
//		return
//	}
//
//	resp, err1 := (&http.Client{}).Do(req)
//	if err1 != nil {
//		response.Data = err1
//		response.Message = "HttpClient Do Failed"
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	//body, err2 := ioutil.ReadAll(resp.Body)
//	//if err2 != nil {
//	//	response.Data = err2
//	//	response.Message = "ReadAll Failed"
//	//	c.JSON(http.StatusOK, response)
//	//	return
//	//}
//	var data struct {
//		Building bool   `json:"building"`
//		Result   string `json:"result"`
//	}
//	if ok := json.NewDecoder(resp.Body).Decode(&data); ok != nil {
//		response.Message = "NewDecoder Failed"
//		response.Data = ok
//		c.JSON(http.StatusOK, response)
//		return
//	}
//	response.Data = data.Result
//
//	//result := service.StatusJob(c.Query("job-name"))
//	//if result == "" {
//	//	response.Message = "Build Ing"
//	//	response.Data = false
//	//} else if result == "ABORTED" {
//	//	response.Message = "Build Aborted"
//	//	response.Data = true
//	//} else if result == "SUCCESS" {
//	//	response.Message = "Build Successful"
//	//	response.Data = true
//	//} else if result == "FAILURE" {
//	//	response.Message = "Build Failed"
//	//	response.Data = true
//	//}
//	c.JSON(http.StatusOK, response)
//}

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
