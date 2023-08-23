package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
)

func IngressList(c *gin.Context) {
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
	ingressList, err := service.IngressList(path, namespace)
	if err != nil {
		response.Data = err
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	if len(ingressList.Items) == 0 {
		response.Message = namespace + " 名称空间中无ingress资源"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = ingressList.Items
	c.JSON(http.StatusOK, response)
}

func IngressListV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}

	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "env namespace 参数不为空"
		c.JSON(http.StatusOK, response)
		return
	}

	ingress := model.Ingress{}
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

func IngressCreate(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 接收参数
	jsonData := model.IngressCreate{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json解析失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据处理
	configFile := jsonData.Env + "config"
	var httpIngressPath []v1beta1.HTTPIngressPath
	for _, data := range jsonData.Rules {
		httpIngressPath = append(httpIngressPath, v1beta1.HTTPIngressPath{
			Path: data.Path,
			PathType: func() *v1beta1.PathType {
				pathType := v1beta1.PathTypePrefix
				return &pathType
			}(),
			Backend: v1beta1.IngressBackend{
				ServiceName: data.TargetService,
				ServicePort: intstr.FromInt(data.TargetPort),
			},
		})
	}
	v1beta1Ingress := &v1beta1.Ingress{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: jsonData.Namespace,
			Name:      jsonData.Name,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "$1",
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: jsonData.Domain,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: httpIngressPath,
						},
					},
				},
			},
		},
	}

	// 创建资源
	result, err := service.IngressCreate(configFile, jsonData.Namespace, v1beta1Ingress)
	if err != nil {
		response.Data = err.Error()
		response.Message = "资源创建失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库操作
	marshalData, err := json.Marshal(jsonData.Rules)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json序列化失败"
		c.JSON(http.StatusOK, response)
		return
	}

	ingress := model.Ingress{}
	ok, err := ingress.Insert(jsonData.Env, jsonData.Namespace, jsonData.Name, jsonData.Domain, string(marshalData))
	if err != nil && !ok {
		response.Message = "数据库操作失败"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = result
	c.JSON(http.StatusOK, response)
}
