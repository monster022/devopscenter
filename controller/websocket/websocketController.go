package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrades = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Demo(c *gin.Context) {
	ws, err := upgrades.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//message = []byte("sss")
		//write ws data
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}

}
