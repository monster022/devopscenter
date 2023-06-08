package middleware

import (
	"devopscenter/configuration"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"log"
	"net/http"
)

func LdapAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Connect to LDAP server

		l, err := ldap.Dial("tcp", configuration.Configs.LdapHost+":"+configuration.Configs.LdapPort)
		if err != nil {
			log.Printf("Error connecting to LDAP server: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer l.Close()

		// Bind with user's credential
		err = l.Bind(fmt.Sprintf("cn=%s,ou=技术部,dc=mojorycorp,dc=cn", username), password)
		if err != nil {
			log.Printf("Error binding with user's credentials: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// User is authenticated, pass control to next handler
		c.Next()
	}
}
