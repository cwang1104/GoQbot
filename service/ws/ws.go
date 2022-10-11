package ws

import (
	"github.com/gorilla/websocket"
	"net/http"
	"qbot/bot"
	"qbot/bot/common/tools"
	"qbot/pkg/logger"
	"time"
)

var wsupgreder = websocket.Upgrader{

	HandshakeTimeout: 5 * time.Second,
	//取消跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func BotWsHandler(w http.ResponseWriter, r *http.Request) {

	var conn *websocket.Conn
	conn, err := wsupgreder.Upgrade(w, r, nil)
	if err != nil {
		logger.Log.Errorf("Upgrade ws failed,err:%v", err)
		return
	}

	logger.Log.Infof("ws connect success,start receiving messages... ")

	//发送消息均为异步发送
	go func() {
		for {
			_ = conn.WriteJSON(<-tools.MsgChan)
		}
	}()

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Errorf("read msg failed,err:%v", err)
			return
		}
		//消息处理
		bot.MessageDistribution(data)
	}
}
