package authority

import (
	"devopscenter/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func List(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}

	name := c.Param("name")
	menu := model.Menu{}
	data, err := menu.ListMenus(name)
	if err != nil {
		response.Code = 50000
		response.Message = err.Error()
	} else {
		response.Data = data
	}

	c.JSON(http.StatusOK, response)
}
