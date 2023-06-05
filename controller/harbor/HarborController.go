package harbor

import (
	"devopscenter/configuration"
	"devopscenter/model"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func List(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	url := configuration.Configs.HarborUrl + "/api/repositories/" + c.Query("env") + "%2F" + c.Query("project") + "/tags"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(configuration.Configs.HarborUsername+":"+configuration.Configs.HarborPassword)))
	if err != nil {
		response.Data = err
		response.Message = "NewRequest Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	resp, err1 := (&http.Client{}).Do(req)
	if err1 != nil {
		response.Data = err1
		response.Message = "HttpClient Do Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		response.Data = err2
		response.Message = "ReadAll Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	var data []*model.HarborTagResponse
	if err3 := json.Unmarshal(body, &data); err3 != nil {
		response.Data = err3
		response.Message = "Unmarshal Failed"
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = data
	c.JSON(http.StatusOK, response)
}
