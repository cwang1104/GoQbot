package socket

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	//消息通道
	news = make(map[string]chan interface{})
	//websocket连接池
	client = make(map[string]*websocket.Conn)
	//互斥锁
	mux sync.Mutex
)

func BotWs(c *gin.Context) {
	BotWsHandler(c.Writer, c.Request)
}

//api接口处理函数
func GetPushNews(c *gin.Context) {
	id := "333"
	wsHandler(c.Writer, c.Request, id)
}

func DeleteClient(c *gin.Context) {
	id := c.Param("id")

	conn, exist := getClient(id)
	if exist {
		conn.Close()
		deleteClient(id)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "未找到客户端",
		})
	}
}

var wsupgreder = websocket.Upgrader{
	//ReadBufferSize:   1024,
	//WriteBufferSize:  1024,
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

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("read msg failed", err)
			return
		}
		log.Println("read msg", string(data))
	}
}

//处理ws请求
func wsHandler(w http.ResponseWriter, r *http.Request, id string) {
	var conn *websocket.Conn
	var err error
	var exist bool
	uid := id
	//创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 2)
	conn, err = wsupgreder.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	_, data, err := conn.ReadMessage()
	if err != nil {
		log.Println("read token err ", err)
		return
	}

	log.Println("msg = ", string(data))

	var msg TokenMsg
	err = json.Unmarshal(data, &msg)
	if err != nil {
		log.Println("unmarshal json err ", err)
		return
	}

	if msg.Token != "token" {
		return
	}
	uid = msg.UserID
	log.Println(uid + "上线")
	//把客户端链接添加到客户端连接池当中
	addClient(uid, conn)

	//获取该客户端的消息通道
	m, exist := getNewsChannel(id)
	if !exist {
		m = make(chan interface{}, 10)
		addNewsChannel(id, m)
	}

	//设置该客户端关闭ws的回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteClient(id)
		log.Println("code", code)
		return nil
	})

	for {
		select {
		case content, _ := <-m:
			//从消息通道接收消息，然后推送给前端
			err := conn.WriteJSON(content)
			if err != nil {
				log.Println(err)
				conn.Close()
				deleteClient(id)
				return
			}
		case <-pingTicker.C:
			//服务端心跳，每20秒ping一次客户端，查看是否在线
			//conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			//err = conn.WriteMessage(websocket.PingMessage, []byte{})
			//if err != nil {
			//	log.Println("send ping err:", err)
			//	conn.Close()
			//	deleteClient(id)
			//	return
			//}

			conn.SetReadDeadline(time.Now().Add(time.Second * 10))
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("read message error", err)
				conn.Close()
				deleteClient(id)
				return
			}
			log.Println(string(data))

			conn.WriteMessage(websocket.BinaryMessage, []byte("pong"))

		}
	}
}

//将客户端链接添加到连接池
func addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	client[id] = conn
	mux.Unlock()
}

//获取指定客户端链接
func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = client[id]
	mux.Unlock()
	return
}

//删除客户端链接
func deleteClient(id string) {
	mux.Lock()
	delete(client, id)
	log.Println(id + "websocket退出")
	mux.Unlock()
}

//添加用户消息通道
func addNewsChannel(id string, m chan interface{}) {
	mux.Lock()
	news[id] = m
	mux.Unlock()
}

//获取指定消息通道
func getNewsChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	m, exist = news[id]
	mux.Unlock()
	return
}

//删除消息通道
func deleteNewsChannel(id string) {
	mux.Lock()
	if m, ok := news[id]; ok {
		close(m)
		delete(news, id)
	}
	mux.Unlock()
}
