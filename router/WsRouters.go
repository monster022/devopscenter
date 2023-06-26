package router

import (
	"devopscenter/controller/bash"
	"devopscenter/controller/websocket"
	"github.com/gin-gonic/gin"
)

func WebsocketRegister(c *gin.Engine) {
	ws := c.Group("/devops")
	{
		ws.GET("/ws", websocket.Demo)
		ws.GET("/:ip/:username/bash/ws", websocket.DemoV2)
		ws.POST("/:ip/:username/bash", bash.Bash)
	}
}
