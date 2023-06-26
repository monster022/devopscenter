package websocket

import (
	"devopscenter/model"
	"devopscenter/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
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
func DemoV2(c *gin.Context) {
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

		operator := model.Operator{}
		operator.Create(c.Query("user"), c.Query("instance"), string(message))
		command := string(message)
		command = strings.Replace(command, "\r", " ", -1)
		machine := model.Machine{}
		password, _ := machine.PasswordByIp(c.Param("ip"))
		output := service.BashCommand(c.Param("ip"), c.Param("username"), password, command)

		//write ws data
		err = ws.WriteMessage(mt, output)
		if err != nil {
			break
		}
	}
}
