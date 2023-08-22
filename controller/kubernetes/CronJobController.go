package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func CronJobList(c *gin.Context) {
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
	cronjob, err := service.CronJobList(path, namespace)
	response.Data = cronjob.Items
	if err != nil {
		response.Data = err
		response.Message = "Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	c.JSON(http.StatusOK, response)
}

func CronJobCreate(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	// 处理请求
	jsonData := model.CronJobCreate{}
	err := c.ShouldBindJSON(&jsonData)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json解析失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据处理
	configFile := jsonData.Env + "config"
	var envVar []coreV1.EnvVar
	for _, data := range jsonData.Data {
		envVar = append(envVar, coreV1.EnvVar{
			Name:  data.Name,
			Value: data.Value,
		})
	}
	v1beta1Cronjob := &v1beta1.CronJob{
		ObjectMeta: metaV1.ObjectMeta{
			Namespace: jsonData.Namespace,
			Name:      jsonData.Name,
		},
		Spec: v1beta1.CronJobSpec{
			Schedule: jsonData.Schedule,
			JobTemplate: v1beta1.JobTemplateSpec{
				Spec: v1.JobSpec{
					Template: coreV1.PodTemplateSpec{
						Spec: coreV1.PodSpec{
							Containers: []coreV1.Container{
								{
									Image:           jsonData.Image,
									ImagePullPolicy: coreV1.PullIfNotPresent,
									Name:            jsonData.Name,
									Env:             envVar,
								},
							},
							RestartPolicy: coreV1.RestartPolicyNever,
						},
					},
				},
			},
		},
	}

	// 创建资源
	job, err := service.CronJobCreate(configFile, jsonData.Namespace, v1beta1Cronjob)
	if err != nil {
		response.Data = err.Error()
		response.Message = "资源创建失败"
		c.JSON(http.StatusOK, response)
		return
	}

	// 数据库操作
	marshalData, err := json.Marshal(jsonData.Data)
	if err != nil {
		response.Data = err.Error()
		response.Message = "json序列化失败"
		c.JSON(http.StatusOK, response)
		return
	}
	cronjob := model.CronJob{}
	_, err = cronjob.Insert(jsonData.Env, jsonData.Namespace, jsonData.Name, jsonData.Schedule, jsonData.Image, string(marshalData))
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = job
	c.JSON(http.StatusOK, response)
}

func CronJobListV2(c *gin.Context) {
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

	// 数据库操作
	cronjob := model.CronJob{}
	data, err := cronjob.List(env, namespace)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}

	response.Data = data
	c.JSON(http.StatusOK, response)
}
