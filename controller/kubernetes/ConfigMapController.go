package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"encoding/json"
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

func ConfigMapListV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 解析参数
	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "env namespace 参数不为空"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库查询
	configmap := model.ConfigMap{}
	data, err := configmap.List(env, namespace)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库查询失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 返回结果
	response.Data = data
	c.JSON(http.StatusOK, response)

}

func ConfigMapAddV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 解析json
	configmapJson := model.ConfigMapJson{}
	err := c.ShouldBindJSON(&configmapJson)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json解析失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// k8s添加
	configFile := configmapJson.Env + "config"
	configmapData := make(map[string]string)
	for _, data := range configmapJson.Data {
		configmapData[data.Key] = data.Value
	}
	configMap := &v1.ConfigMap{}
	configMap.Namespace = configmapJson.Namespace
	configMap.Data = configmapData
	configMap.Name = configmapJson.Name

	result, err := service.ConfigMapAdd(configFile, configmapJson.Namespace, configMap)
	if err != nil {
		response.Data = err.Error()
		response.Message = "K8s添加configMap: " + configmapJson.Name + " 失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库添加数据
	config := model.ConfigMap{}
	baseData := model.ConfigBase{
		Env:       configmapJson.Env,
		Name:      configmapJson.Name,
		Namespace: configmapJson.Namespace,
	}
	marshalData, err := json.Marshal(configmapJson.Data)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json序列化失败"
		c.JSON(http.StatusOK, response)
		return
	}

	ok, err := config.Insert(baseData, string(marshalData))
	if err != nil && !ok {
		response.Data = err.Error()
		response.Message = "数据库插入失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = result
	c.JSON(http.StatusOK, response)
}
