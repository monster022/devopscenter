package order

import (
	"devopscenter/model"
	"devopscenter/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func List(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "success",
		Data:    nil,
	}
	pageString := c.Query("page")
	sizeString := c.Query("size")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	order := model.Order{}
	data, err := order.ListOrder(page, size)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Database Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func ListTackleName(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "success",
		Data:    nil,
	}
	pageString := c.Query("page")
	sizeString := c.Query("size")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	order := model.Order{}
	data, err := order.ListTackleName(page, size, c.Param("name"))
	if err != nil {
		response.Data = err.Error()
		response.Message = "Database Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func ListTackleNameTotal(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "success",
		Data:    nil,
	}
	order := model.Order{}
	data, err := order.ListTackleNameCount(c.Param("name"))
	if err != nil {
		response.Data = err.Error()
		response.Message = "Databases Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func ListSubmitName(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "success",
		Data:    nil,
	}
	pageString := c.Query("page")
	sizeString := c.Query("size")
	page, err := strconv.Atoi(pageString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	size, err := strconv.Atoi(sizeString)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	order := model.Order{}
	data, err := order.ListSubmitName(page, size, c.Param("name"))
	if err != nil {
		response.Data = err.Error()
		response.Message = "Database Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}

func PatchOrder(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	order := model.Order{}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	_, err = order.PatchOrderStatus(id, c.Param("status"))
	if err != nil {
		response.Data = err.Error()
		response.Message = "Sql Exec Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	//if ok := utils.Alter(order.NameByOrderId(id), "您的待办事项已处理，请及时查看"); ok == 0 {
	//	response.Message = "FeishuTalk Notify Failed"
	//	c.JSON(http.StatusOK, response)
	//	return
	//}
	response.Data = true
	c.JSON(http.StatusOK, response)
}

func PostOrder(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	order := model.Order{}
	var r struct {
		RejectReason string `json:"rejectReason"`
	}
	if err := c.ShouldBindJSON(&r); err != nil {
		response.Message = "Paras Json Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	_, err = order.PostOrderStatusReject(id, c.Param("status"), r.RejectReason)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Sql Exec Failed"
		c.JSON(http.StatusOK, response)
		return
	}

	//if ok := utils.Alter(order.NameByOrderId(id), "您的待办事项已处理，请及时查看"); ok == 0 {
	//	response.Message = "FeishuTalk Notify Failed"
	//	c.JSON(http.StatusOK, response)
	//	return
	//}

	response.Data = true
	c.JSON(http.StatusOK, response)
}

func Create(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	order := model.Order{}
	orderRequestBody := model.OrderRequestBody{}
	if err := c.ShouldBindJSON(&orderRequestBody); err != nil {
		response.Data = err.Error()
		response.Message = "Paras Json Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	_, err := order.CreateOrder(&orderRequestBody)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Sql Exec Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	if ok := utils.Alter(orderRequestBody.TackleName, "您有新的待办事项，请及时处理"); ok == 0 {
		response.Message = "FeishuTalk Notify Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
}
