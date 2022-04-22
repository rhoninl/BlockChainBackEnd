package Controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"main/Model"
	"net/http"
)

func BuildSocket(c *gin.Context) {
	id, _ := c.Get("companyId")
	upGrader := websocket.Upgrader{ // 将http协议升级成websocket
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Subprotocols: []string{c.GetHeader("Sec-WebSocket-Protocol")},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil) // 获取链接
	if err != nil {                                         // 错误检测
		log.Println("websocket connect default ", err, c.Param("channel"))
		c.JSON(http.StatusBadRequest, gin.H{"message": "Websocket建立错误"})
		return
	}
	Model.UseClient().Login(conn, id.(int64)) // 将链接信息交到管理中心
}
