package helper

import (
	"devopscenter/configuration"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"log"
)

func OpenldapVerify(username, password string) bool {
	l, err := ldap.Dial("tcp", configuration.Configs.LdapHost+":"+configuration.Configs.LdapPort)
	if err != nil {
		log.Printf("Error connecting to LDAP server: %v", err)
		return false
	}
	defer l.Close()
	err = l.Bind(fmt.Sprintf("cn=%s,ou=技术部,dc=mojorycorp,dc=cn", username), password)
	if err != nil {
		log.Printf("Error binding with user's credentials: %v", err)
		return false
	}
	return true
}

func OpenldapModifyPassword(username, newPassword string) bool {
	l, err := ldap.Dial("tcp", configuration.Configs.LdapHost+":"+configuration.Configs.LdapPort)
	if err != nil {
		log.Printf("Error connecting to LDAP server: %v", err)
		return false
	}
	defer l.Close()

	err = l.Bind("cn=admin,dc=mojorycorp,dc=cn", "mojory@1q2w3e4r")
	if err != nil {
		return false
	}
	modifyRequest := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=技术部,dc=mojorycorp,dc=cn", username), []ldap.Control{})
	modifyRequest.Replace("userPassword", []string{newPassword})

	// 执行密码修改请求
	err = l.Modify(modifyRequest)
	if err != nil {
		return false
	}
	return true
}
