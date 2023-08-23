package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"net/http"
	"strconv"
	"strings"
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

	// 记录数据库发布的版本
	deployInfo := model.DeployProjectDetail{}
	commitId := strings.Split(data.ImageSource, "-")
	_, err := deployInfo.CreateDeployInfo(data.DeploymentName, commitId[1], data.Env, data.CreateBy, data.Namespace, data.PublishType, data.ImageSource)
	if err != nil {
		response.Message = "发布历史记录数据库失败"
		response.Data = err.Error()
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

func DeployListV3(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 参数处理
	env := c.Query("env")
	namespace := c.Query("namespace")
	if env == "" || namespace == "" {
		response.Message = "env namespace 参数不能为空"
		c.JSON(http.StatusOK, response)
		return
	}

	deploy := model.DeployAdd{}
	data, err := deploy.List(env, namespace)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = data
	c.JSON(http.StatusOK, response)
}

func DeployAddV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 解析json
	jsonData := model.DeployAdd{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json解析失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据处理
	replicas := int32(jsonData.Replicas)
	configFile := jsonData.Env + "config"
	deployAdd := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Labels: map[string]string{
				"app": jsonData.Name,
			},
			Name:      jsonData.Name,
			Namespace: jsonData.Namespace,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": jsonData.Name,
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{
						"app": jsonData.Name,
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							EnvFrom: []coreV1.EnvFromSource{
								{
									ConfigMapRef: &coreV1.ConfigMapEnvSource{
										LocalObjectReference: coreV1.LocalObjectReference{
											Name: jsonData.ConfigName,
										},
										Optional: boolPointer(false),
									},
								},
							},
							Image:           jsonData.Image,
							ImagePullPolicy: coreV1.PullIfNotPresent,
							Name:            jsonData.Name,
							LivenessProbe: &coreV1.Probe{
								FailureThreshold: 3,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								TimeoutSeconds:   1,
								ProbeHandler: coreV1.ProbeHandler{
									HTTPGet: &coreV1.HTTPGetAction{
										Path:   jsonData.Uri,
										Port:   intstr.IntOrString{Type: intstr.Int, IntVal: jsonData.Port},
										Scheme: "HTTP",
									},
								},
							},
							ReadinessProbe: &coreV1.Probe{
								FailureThreshold: 3,
								PeriodSeconds:    10,
								SuccessThreshold: 1,
								TimeoutSeconds:   1,
								ProbeHandler: coreV1.ProbeHandler{
									HTTPGet: &coreV1.HTTPGetAction{
										Path:   jsonData.Uri,
										Port:   intstr.IntOrString{Type: intstr.Int, IntVal: jsonData.Port},
										Scheme: "HTTP",
									},
								},
							},
							Resources: coreV1.ResourceRequirements{
								Requests: coreV1.ResourceList{
									coreV1.ResourceCPU:    resource.MustParse(jsonData.Cpu),
									coreV1.ResourceMemory: resource.MustParse(jsonData.Mem),
								},
								Limits: coreV1.ResourceList{
									coreV1.ResourceCPU:    resource.MustParse(jsonData.Cpu),
									coreV1.ResourceMemory: resource.MustParse(jsonData.Mem),
								},
							},
						},
					},
				},
			},
		},
	}

	// 创建deployment
	result, err := service.DeploymentAdd(configFile, jsonData.Namespace, deployAdd)
	if err != nil {
		response.Data = err.Error()
		response.Message = "创建Deployment失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库记录
	ok, err := jsonData.Insert(&jsonData)
	if err != nil && !ok {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = result
	c.JSON(http.StatusOK, response)
}

func boolPointer(b bool) *bool {
	return &b
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

func DeployDelete(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 数据处理
	env := c.Query("env")
	namespace := c.Query("namespace")
	deployment := c.Query("deployment")
	idParam := c.Query("id")
	if env == "" || namespace == "" || deployment == "" {
		response.Message = "env namespace deployment 参数不能为空"
		c.JSON(http.StatusOK, response)
		return
	}

	// 删除资源
	configFile := env + "config"
	err := service.DeploymentDelete(configFile, deployment, namespace)
	if err != nil {
		response.Data = err.Error()
		response.Message = "资源删除失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 删除数据库记录
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Data = err.Error()
		response.Message = "类型转换失败"
		c.JSON(http.StatusOK, response)
		return
	}
	deploy := model.DeployAdd{}
	ok, err := deploy.Delete(id)
	if err != nil && !ok {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
}
