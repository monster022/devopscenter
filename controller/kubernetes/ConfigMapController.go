package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

func ConfigMapList(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	configFile := c.Query("env") + "config"
	data, err := service.ConfigMapList(configFile, c.Query("namespace"))
	if err != nil {
		response.Message = "该环境或者名称空间中无资源"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data.Items
	c.JSON(http.StatusOK, response)
}

func ConfigMapAdd(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	json := model.ConfigMap{}
	c.ShouldBindJSON(&json)
	configFile := json.Env + "config"
	data := make(map[string]string)
	data["ASPNETCORE_ENVIRONMENT"] = json.Env
	configMap := &v1.ConfigMap{}
	configMap.Namespace = json.Namespace
	configMap.Data = data
	configMap.Name = json.Name
	result, err := service.ConfigMapAdd(configFile, json.Namespace, configMap)
	if err != nil {
		response.Message = "failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = result
	c.JSON(http.StatusOK, response)
}
