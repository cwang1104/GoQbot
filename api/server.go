package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"qbot/api/http"
	"qbot/api/middleware"
	"qbot/api/socket"
)

type Server struct {
	Port   string
	Router *gin.Engine
}

func NewServer(port string) *Server {
	server := Server{
		Port: port,
	}
	server.setRouter()
	return &server
}

func (s *Server) setRouter() {
	r := gin.Default()

	//跨域请求
	r.Use(middleware.Cors())

	//注册与登录
	r.POST("/register", http.UserRegister)
	r.POST("/login", http.CheckLogin)

	//机器人操作
	botGroup := r.Group("/bot")
	{
		//token校验
		botGroup.Use(middleware.AuthToken())

		//  @全体成员
		botGroup.POST("/at_all_member", http.AtAllMember)

		//定时器测试
		botGroup.POST("/cro", http.AddDsq)
		botGroup.POST("/del", http.DelJob)
		botGroup.POST("/list", http.JobList)
	}

	//websocket服务
	r.GET("/ws/bot", socket.BotWs)
	s.Router = r
}

func (s *Server) RunServer(address string) error {
	url := fmt.Sprintf("%s:%s", address, s.Port)
	return s.Router.Run(url)
}
