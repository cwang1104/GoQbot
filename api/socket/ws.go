package socket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"qbot/bot/ws"
	"time"
)

func BotWs(c *gin.Context) {
	BotWsHandler(c.Writer, c.Request)
}

var wsupgreder = websocket.Upgrader{

	HandshakeTimeout: 5 * time.Second,
	//取消跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type TokenMsg struct {
	Code   int    `json:"code"`
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func BotWsHandler(w http.ResponseWriter, r *http.Request) {

	var conn *websocket.Conn
	conn, err := wsupgreder.Upgrade(w, r, nil)
	if err != nil {
		log.Println("ws升级失败", err)
		return
	}

	log.Println("接收开始")

	//发送消息均为异步发送
	go func() {
		for {
			_ = conn.WriteJSON(<-ws.MsgChan)
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read msg failed", err)
			return
		}

		//将上报的消息反序列化为结构体
		var messageMode ws.MessageType
		_ = json.Unmarshal(data, &messageMode)

		//上报为消息时传入消息处理
		if messageMode.PostType == "message" {
			ws.MessageDistribution(&messageMode)
		}
	}
}
