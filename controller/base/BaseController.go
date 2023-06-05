package base

import (
	"devopscenter/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"language": "go",
		"type":     "application callback",
	})
}

func PostJson(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	// 解析json
	user := model.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Message = "json 解析失败"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	response.Message = "json 解析成功"
	response.Data = user
	c.JSON(http.StatusOK, response)
}

func PostForm(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "successful",
		Data:    nil,
	}
	name, _ := c.GetPostForm("name")
	age, _ := c.GetPostForm("age")
	phone, _ := c.GetPostForm("phone")
	data := make(map[string]string)
	data["name"] = name
	data["age"] = age
	data["phone"] = phone
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func Delete(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "Delete Successful",
		Data:    nil,
	}
	name := c.Query("name")
	if name == "" {
		response.Message = "Delete Fail"
		response.Data = name
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = name
	c.JSON(http.StatusOK, response)
}

func Patch(c *gin.Context) {
	response := model.Res{
		Code:    200,
		Message: "Patch Successful",
		Data:    nil,
	}
	name := c.Query("name")
	if name == "" {
		response.Message = "Patch Fail"
		response.Data = name
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = name
	c.JSON(http.StatusOK, response)
}
