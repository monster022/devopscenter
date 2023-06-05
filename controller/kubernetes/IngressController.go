package kubernetes

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	ingress := model.Ingress{}
	data := ingress.List(env, namespace, page, size)
	total := ingress.Count(namespace)
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    data,
		"total":   total,
	})
}

func IngressDelete(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	env := c.Query("env")
	namespace := c.Query("namespace")
	ingressName := c.Query("name")
	idParams := c.Query("id")
	if env == "" || namespace == "" || ingressName == "" || idParams == "" {
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
	configFile := env + "config"
	if err := service.IngressDelete(configFile, namespace, ingressName); err != nil {
		response.Message = "Ingress Delete Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	ingress := model.Ingress{}
	if err := ingress.Delete(id); err == false {
		response.Message = "Databases Delete Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
	return
}
