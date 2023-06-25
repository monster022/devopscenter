package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"strconv"
)

func DeployList(c *gin.Context) {
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
	deploymentList, err := service.DeploymentList(path, namespace)
	if err != nil {
		response.Data = err
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	if len(deploymentList.Items) == 0 {
		response.Message = namespace + " 名称空间中无deployment资源"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = deploymentList.Items
	c.JSON(http.StatusOK, response)
}

func DeployGet(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	env := c.Query("env")
	namespace := c.Query("namespace")
	deploymentName := c.Param("name")
	if env == "" || namespace == "" || deploymentName == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	path := env + "config"
	data, err := service.DeploymentGet(path, namespace, deploymentName)
	if err != nil {
		response.Message = "获取失败"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data.Spec.Template.Spec.Containers
	c.JSON(http.StatusOK, response)
}

func DeployListV2(c *gin.Context) {
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
	deployment := model.Deployment{}
	data := deployment.List(env, namespace, page, size)
	total := deployment.Count(namespace)
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    data,
		"total":   total,
	})
}

func DeployPatch(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	data := model.DeploymentImage{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Message = "Json Paras Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	path := data.Env + "config"
	if err := service.DeploymentImagePatch(path, data.Namespace, data.ContainerName, data.DeploymentName, data.ImageSource); err != nil {
		response.Message = "Service Patch Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	deploy := model.Deployment{}
	number := deploy.ImagePatch(data.ImageSource, data.Env, data.Namespace, data.DeploymentName)
	response.Data = number
	if number == -1 {
		response.Message = "Databases Execute Failed"
		c.JSON(http.StatusOK, response)
		return
	} else if number == 0 {
		response.Message = "Databases Not Modify"
		c.JSON(http.StatusOK, response)
		return
	} else if number > 1 {
		response.Message = "Modify Lot Rows"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func DeployAdd(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	json := model.DeploymentAdd{}
	c.ShouldBindJSON(&json)
	replicas := int32(json.Replicas)
	configFile := json.Env + "config"
	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Labels: map[string]string{
				"app": json.Name,
			},
			Name:      json.Name,
			Namespace: json.Namespace,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": json.Name,
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": json.Name,
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name:            json.Name,
							Image:           json.Image,
							ImagePullPolicy: coreV1.PullIfNotPresent,
							EnvFrom: []coreV1.EnvFromSource{
								{
									ConfigMapRef: &coreV1.ConfigMapEnvSource{
										LocalObjectReference: coreV1.LocalObjectReference{
											Name: json.Name,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	data, err := service.DeploymentAdd(configFile, json.Namespace, deployment)
	response.Data = data
	if err != nil {
		response.Message = "deployment create failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func PodList(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	podList, err := service.PodList(c.Param("env")+"config", c.Param("namespace"))
	if err != nil {
		response.Message = "Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = podList
	c.JSON(http.StatusOK, response)
}
