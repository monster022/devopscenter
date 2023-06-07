package middleware

import (
	"devopscenter/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

var JwtKey = []byte("my_secret_key")

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := model.Res{
			Code:    20000,
			Message: "successful",
			Data:    nil,
		}

		tokenString := c.Request.Header.Get("Authorization")

		if tokenString == "" {
			response.Message = "Authorization header required"
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return JwtKey, nil
		})

		if err != nil {
			response.Message = "Invalid token"
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*model.JwtClaims); ok && token.Valid {
			c.Set("username", claims.Username)
			c.Next()
		} else {
			response.Message = "Invalid token"
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
	}
}
