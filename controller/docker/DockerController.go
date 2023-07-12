package docker

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
)

const (
	HttpNewRequestError = 1000
	HttpClientDoError   = 1001
)

func Create(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	data := model.DockerRequestBody{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Message = "Json解析失败"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	// 镜像名称中带有 / 处理
	image := strings.Split(data.Image, ":")

	// 下载镜像
	var wgImagePull sync.WaitGroup
	for _, m := range data.Target {
		wgImagePull.Add(1)
		go service.ImagePull(m, image[0], image[1], &wgImagePull)
	}
	wgImagePull.Wait()

	// 删除镜像
	var wgContainerDelete sync.WaitGroup
	for _, m := range data.Target {
		wgContainerDelete.Add(1)
		go service.ContainerDelete(m, data.Name, &wgContainerDelete)
	}
	wgContainerDelete.Wait()

	// 环境变量处理
	aspEnvironment := data.Env
	switch data.Env {
	case "uat":
		aspEnvironment = "UAT"
		break
	case "fat":
		aspEnvironment = "FAT"
		break
	case "pro":
		aspEnvironment = "Production"
		break
	}

	// 创建镜像
	var wgContainerCreate sync.WaitGroup
	for _, m := range data.Target {
		wgContainerCreate.Add(1)
		go service.ContainerCreate(m, aspEnvironment, data.Image, data.Name, data.Port, &wgContainerCreate)
	}
	wgContainerCreate.Wait()

	// 启动镜像
	var wgContainerStart sync.WaitGroup
	for _, m := range data.Target {
		wgContainerStart.Add(1)
		go service.ContainerStart(m, data.Name, &wgContainerStart)
	}
	wgContainerStart.Wait()

	response.Data = true
	c.JSON(http.StatusOK, response)
}
