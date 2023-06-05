package machine

import (
	"devopscenter/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}
	pageParameter := c.Query("page")
	sizeParameter := c.Query("size")
	if pageParameter == "" || sizeParameter == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	page, err1 := strconv.Atoi(pageParameter)
	size, err2 := strconv.Atoi(sizeParameter)
	if err1 != nil || err2 != nil {
		response.Message = "Type Converter Fail"
		c.JSON(http.StatusOK, response)
		return
	}
	machine := model.Machine{}
	total := machine.Total()
	response.Data = machine.List(page, size)
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    response.Data,
		"total":   total,
	})
}

func Password(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}
	idParameter := c.Query("id")
	if idParameter == "" {
		response.Message = "Parameter Not Found"
		c.JSON(http.StatusOK, response)
		return
	}
	id, err := strconv.Atoi(idParameter)
	if err != nil {
		response.Data = err
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	machine := model.Machine{}
	response.Data = machine.PasswordList(id)
	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	response := model.Res{
		Code:    201,
		Message: "Successful",
		Data:    nil,
	}
	machine := model.Machine{}
	if err := c.ShouldBindJSON(&machine); err != nil {
		response.Message = "json 解析失败"
		response.Data = err
		c.JSON(http.StatusCreated, response)
		return
	}
	if result := machine.Insert(); result == false {
		response.Message = "Databases Insert Failed"
		response.Data = result
		c.JSON(http.StatusCreated, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusCreated, response)
}

func Remove(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}
	idParameter := c.Query("id")
	if idParameter == "" {
		response.Message = "Parameter nil"
		c.JSON(http.StatusOK, response)
		return
	}
	id, err := strconv.Atoi(idParameter)
	if err != nil {
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	machine := model.Machine{}
	if result := machine.Delete(id); result == false {
		response.Message = "Sql Exec Failed"
		response.Data = result
		c.JSON(http.StatusOK, response)
	}
	c.JSON(http.StatusOK, response)
}
func Update(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	machine := model.Machine{}
	if err := c.ShouldBindJSON(&machine); err != nil {
		response.Message = "json 解析失败"
		response.Data = err
		c.JSON(http.StatusCreated, response)
		return
	}
	result, err1 := machine.Update()
	if err1 != nil {
		response.Data = result
		response.Message = "Databases Exec Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = result
	c.JSON(http.StatusCreated, response)
}
