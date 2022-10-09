package socket

import (
	"github.com/gin-gonic/gin"
	"qbot/service/ws"
)

//BotWs 请求升级为websocket请求
//监听go-cqhttp客户端上报的消息，再分发进行处理
func BotWs(c *gin.Context) {
	ws.BotWsHandler(c.Writer, c.Request)
}
