package machine

import (
	"bytes"
	"devopscenter/model"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"net/http"
	"strconv"
)

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
		Code:    20000,
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

func PatchName(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	machine := model.Machine{}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Data = err
		response.Message = "Type Convert Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	if result, ok := machine.PatchName(id, c.Query("name")); ok != nil {
		response.Data = result
		response.Message = "Databases Exec Failed"
		c.JSON(http.StatusOK, response)
		return
	} else {
		response.Data = result
		c.JSON(http.StatusOK, response)
	}
}

func ListV2(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "SUCCESS",
		Data:    nil,
	}
	pageParameter := c.Query("page")
	sizeParameter := c.Query("size")
	ip := c.Query("ip")
	if pageParameter == "" {
		response.Message = "page 参数不能为空"
		response.Data = false
		c.JSON(http.StatusOK, response)
		return
	}
	if sizeParameter == "" {
		response.Message = "size 参数不能为空"
		response.Data = false
		c.JSON(http.StatusOK, response)
		return
	}
	page, err := strconv.Atoi(pageParameter)
	if err != nil {
		response.Data = err.Error()
		response.Message = "类型转换失败"
		c.JSON(http.StatusOK, response)
		return
	}
	size, err := strconv.Atoi(sizeParameter)
	if err != nil {
		response.Data = err.Error()
		response.Message = "类型转换失败"
		c.JSON(http.StatusOK, response)
		return
	}
	machine := model.Machine{}
	data, err := machine.VagueSearch(ip, page, size)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}
	total, err := machine.VagueSearchTotal(ip)
	if err != nil {
		response.Data = err.Error()
		response.Message = "数据库操作失败"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, gin.H{
		"code":    response.Code,
		"message": response.Message,
		"data":    response.Data,
		"total":   total,
	})
}

func Download(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}

	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")

	// Add headers
	row := sheet.AddRow()
	row.AddCell().Value = "序号"
	row.AddCell().Value = "实例名称"
	row.AddCell().Value = "实例地址"
	row.AddCell().Value = "用户名"
	row.AddCell().Value = "密码"
	row.AddCell().Value = "CPU"
	row.AddCell().Value = "内存"
	row.AddCell().Value = "实例状态"
	row.AddCell().Value = "实例标签"

	machine := model.Machine{}
	downloadData, err := machine.DownloadData()
	if err != nil {
		response.Data = err.Error()
		response.Message = "查询数据失败"
		c.JSON(http.StatusOK, response)
		return
	}

	var data [][]string
	for _, valuePtr := range downloadData {
		value := *valuePtr
		row := []string{strconv.Itoa(value.Id), value.InstanceName, value.InstanceIp, value.InstanceUsername, value.InstancePassword,
			strconv.Itoa(value.InstanceCpu), strconv.Itoa(value.InstanceMemory), value.InstanceStatus, value.InstanceTag}
		data = append(data, row)
	}

	for _, rowData := range data {
		row := sheet.AddRow()
		for _, cellValue := range rowData {
			cell := row.AddCell()
			cell.Value = cellValue
		}
	}

	// Save the Excel file to a buffer
	buffer := new(bytes.Buffer)
	err = file.Write(buffer)
	if err != nil {
		response.Data = err.Error()
		response.Message = "文件创建失败"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Set response headers
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=machine.xlsx")

	// Write the buffer to the response writer
	_, err = buffer.WriteTo(c.Writer)
	if err != nil {
		response.Data = err.Error()
		response.Message = "Error writing to response"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
}

func Export(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Successful",
		Data:    nil,
	}
	machine := model.Machine{}
	data, err := machine.DownloadData()
	if err != nil {
		response.Data = err.Error()
		response.Message = "查询数据失败"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}
