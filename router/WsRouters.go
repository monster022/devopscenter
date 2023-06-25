package router

import (
	"devopscenter/controller/websocket"
	"github.com/gin-gonic/gin"
)

func WebsocketRegister(c *gin.Engine) {
	api := c.Group("/devops")
	{
		api.GET("/ws", websocket.Demo)
	}
}
