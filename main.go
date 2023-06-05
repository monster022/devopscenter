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
	// 配置 Cors 中间件
	//config := cors.DefaultConfig()
	//config.AllowOrigins = []string{"http://localhost:3000"} // 允许的前端域名和端口
	//config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	//app.Use(cors.New(config))

	// 日志中间件
	app.Use(middleware.Log())

	router.BaseRegister(app)
	router.MachineRegister(app)
	router.GitlabRegister(app)
	router.JenkinsRegister(app)
	router.HarborRegister(app)
	router.KubernetesRegister(app)
	app.Run("0.0.0.0:" + configuration.Configs.ServerPort)
}
