package main

import (
	"devopscenter/configuration"
	"devopscenter/helper"
	"devopscenter/middleware"
	"devopscenter/router"
	"github.com/gin-gonic/gin"
)

func init() {
	configuration.LoadConfig("./config.ini")
	helper.MysqlConnect()
	helper.GitlabConnect()
	helper.JenkinsConnect()
}

func main() {
	app := gin.Default()

	//跨域中间件
	app.Use(middleware.Cors())

	//登录中间件
	//app.Use(middleware.LdapAuth())

	// 日志中间件
	app.Use(middleware.Log())

	router.BaseRegister(app)
	router.MachineRegister(app)
	router.OrderRegister(app)
	router.GitlabRegister(app)
	router.JenkinsRegister(app)
	router.HarborRegister(app)
	router.KubernetesRegister(app)
	app.Run("0.0.0.0:" + configuration.Configs.ServerPort)
}
