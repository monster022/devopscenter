package login

import (
	"devopscenter/helper"
	"devopscenter/middleware"
	"devopscenter/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Auth(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "Login Successful",
		Data:    nil,
	}
	var loginValues struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginValues); err != nil {
		response.Data = err
		response.Message = "Invalid request"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if ok := helper.OpenldapVerify(loginValues.Username, loginValues.Password); ok == false {
		response.Code = 50001
		response.Message = "username or password error"
		response.Data = ok
		c.JSON(http.StatusOK, response)
		return
	}
	expirationTime := time.Now().Add(6 * time.Hour)
	claims := &model.JwtClaims{
		Username: loginValues.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtKey)
	if err != nil {
		response.Message = "Could not generate token"
		response.Data = err
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	type temp struct {
		Token string `json:"token"`
	}
	response.Data = temp{tokenString}
	c.JSON(http.StatusOK, response)
}

func ModifyPassword(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	var entry struct {
		Username        string `json:"username"`
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	if err := c.ShouldBindJSON(&entry); err != nil {
		response.Message = "Json Paras Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
	}
	if err := helper.OpenldapVerify(entry.Username, entry.CurrentPassword); err == false {
		response.Message = "Current Password Error"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	if err := helper.OpenldapModifyPassword(entry.Username, entry.NewPassword); err != true {
		response.Message = "Modify Password Failed"
		response.Data = err
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = true
	c.JSON(http.StatusOK, response)
}
