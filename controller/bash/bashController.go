package bash

import (
	"devopscenter/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
	"net/http"
)

func Bash(c *gin.Context) {
	response := model.Res{
		Code:    20000,
		Message: "successful",
		Data:    nil,
	}
	machine := model.Machine{}
	password, err := machine.PasswordByIp(c.Param("ip"))
	if err != nil {
		response.Message = "Machine Password Not Found"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}

	config := &ssh.ClientConfig{
		User: c.Param("username"),
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", c.Param("ip")+":22", config)
	if err != nil {
		response.Message = "Cannot Connect " + c.Param("ip")
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	defer client.Close()

	// 执行命令
	session, err := client.NewSession()
	if err != nil {
		response.Message = "Create Session Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	defer session.Close()

	output, err := session.CombinedOutput("ls -l")
	if err != nil {
		response.Message = "Execute ls -l Command Failed"
		response.Data = err.Error()
		c.JSON(http.StatusOK, response)
		return
	}
	response.Data = string(output)
	c.JSON(http.StatusOK, response)
}
