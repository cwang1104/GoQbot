package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"qbot/api/bothttp"
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

	//  @全体成员
	r.GET("/at_all_member", bothttp.AtAllMember)

	//定时器测试
	r.GET("/cro", bothttp.AddDsq)
	r.GET("/del", bothttp.DelJob)
	r.GET("/list", bothttp.JobList)

	//websocket服务
	r.GET("/ws/bot", socket.BotWs)
	s.Router = r
}

func (s *Server) RunServer(address string) error {
	url := fmt.Sprintf("%s:%s", address, s.Port)
	return s.Router.Run(url)
}
