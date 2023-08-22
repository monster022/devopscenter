package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	V1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
)

func ServiceList(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	path := env + "config"
	if namespace == "" {
		namespace = "default"
	}
	serviceList, err := service.SvcList(path, namespace)
	if err != nil {
		response.Data = err
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	if len(serviceList.Items) == 0 {
		response.Message = namespace + " 名称空间中无service资源"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = serviceList.Items
	c.JSON(http.StatusOK, response)
}

func ServiceListV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "env namespace 参数不为空"
		c.JSON(http.StatusOK, response)
		return
	}

	ingress := model.Service{}
	data, err := ingress.List(env, namespace)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = data
	c.JSON(http.StatusOK, response)
}

func ServiceCreateV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}

	// 解析数据
	jsonData := model.ServiceCreateV2{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json解析失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据处理
	configFile := jsonData.Env + "config"
	svc := &V1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: jsonData.Namespace,
			Name:      jsonData.Name,
		},
		Spec: V1.ServiceSpec{
			Ports: []V1.ServicePort{
				0: {
					Name:       jsonData.Name,
					Port:       int32(jsonData.Port),
					Protocol:   V1.Protocol(jsonData.Protocol),
					TargetPort: intstr.FromInt(jsonData.TargetPort),
				},
			},
			Selector: map[string]string{
				"app": jsonData.Deployment,
			},
			Type: V1.ServiceType(jsonData.Type),
		},
	}

	// 创建Service
	result, err := service.SvcCreate(configFile, jsonData.Namespace, svc)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Service: " + jsonData.Name + " 创建失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库记录
	service := model.Service{}
	ok, err := service.Insert(&jsonData)
	if err != nil && !ok {
		response.Data = err.Error()
		response.Message = "数据库执行失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = result
	c.JSON(http.StatusOK, response)
}
