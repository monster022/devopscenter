package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	V1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
	"strconv"
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
		Message: "successful",
		Data:    nil,
	}
	projectPage := c.Query("page")
	projectSize := c.Query("size")
	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	page, err1 := strconv.Atoi(projectPage)
	size, err2 := strconv.Atoi(projectSize)
	if err1 != nil || err2 != nil {
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	service := model.Service{}
	data := service.List(env, namespace, page, size)
	total := service.Count(namespace)
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    data,
		"total":   total,
	})
}

func ServiceDelete(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	env := c.Query("env")
	namespace := c.Query("namespace")
	serviceName := c.Query("name")
	idParams := c.Query("id")
	if env == "" || namespace == "" || serviceName == "" || idParams == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	id, err := strconv.Atoi(idParams)
	if err != nil {
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	path := env + "config"
	if err := service.SvcDelete(path, namespace, serviceName); err != nil {
		response.Message = "Service Delete Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	service := model.Service{}
	if err := service.Delete(id); err == false {
		response.Message = "Databases Delete Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
}

func ServiceCreate(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	env := c.Query("env")
	namespace := c.Query("namespace")
	json := model.ServiceCreate{}
	if err := c.ShouldBindJSON(&json); err != nil {
		response.Message = "Json ShouldBindJSON Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	if env == "" || namespace == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	svc := &V1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: namespace,
			Name:      json.Name,
		},
		Spec: V1.ServiceSpec{
			Ports: []V1.ServicePort{
				0: {
					/*
						Name:       s.Name,
						Protocol:   apiv1.ProtocolTCP,
						Port:       utils.StrToInt32(s.Port),
						TargetPort: intstr.IntOrString{intstr.Int, utils.StrToInt32(s.Port), s.Port},
					*/
					Name:     json.Name,
					Protocol: V1.Protocol(json.Protocol),
					Port:     int32(json.Port),
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: int32(json.TargetPort),
						StrVal: "TCP",
					},
				},
			},
			Type: V1.ServiceType(json.Type),
		},
	}
	configFile := env + "config"
	data, err := service.SvcCreate(configFile, namespace, svc)
	if err != nil {
		response.Message = "Service Create Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	service := model.Service{}
	service.Name = json.Name
	service.PortName = json.Name
	service.Port = json.Port
	service.TargetPort = strconv.Itoa(json.TargetPort)
	service.Protocol = json.Protocol
	service.Type = json.Type
	service.Env = env
	service.Namespace = namespace
	if json.Type == "NodePort" {
		service.NodePort = int(data.Spec.Ports[0].NodePort)
	} else {
		service.NodePort = 0
	}
	if err := service.Insert(); err == false {
		response.Message = "Databases Insert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}
