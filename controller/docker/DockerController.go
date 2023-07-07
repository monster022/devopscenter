package docker

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	"net/http"
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

	// 定义一个channel
	statusChan := make(chan int)

	// 删除容器
	for _, machine := range data.Target {
		go service.DeleteContainer(machine, data.Name, statusChan)
	}
	// 通过 channel 获取值来判断容器是否删除成功
	for i := 0; i < len(data.Target); i++ {
		s := <-statusChan
		if s == 2 {
			response.Message = "删除容器失败"
			response.Data = s
			c.JSON(http.StatusOK, response)
			return
		}
	}

	// 创建容器
	for _, machine := range data.Target {
		go service.CreateContainer(machine, data.Name, data.Env, data.Image, data.Port, statusChan)
	}
	// 通过 channel 获取值来判断你是否创建成功， 成功后则启动容器
	for _, machine := range data.Target {
		s := <-statusChan
		if s == 0 {
			response.Message = "容器创建失败"
			response.Data = false
			c.JSON(http.StatusOK, response)
			return
		} else {
			go service.StartContainer(machine, data.Name, statusChan)
		}
	}
	// 通过 channel 获取值来判断容器是否启动成功
	for i := 0; i < len(data.Target); i++ {
		s := <-statusChan
		if s == 1 {
			response.Message = "容器启动失败"
			response.Data = false
			c.JSON(http.StatusOK, response)
			return
		}
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
}
