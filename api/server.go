package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"qbot/api/http"
	"qbot/api/socket"
	"qbot/middleware"
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

	r.GET("/map", http.ShowCronMap)

	//机器人操作
	botGroup := r.Group("/bot")
	{
		//token校验
		botGroup.Use(middleware.AuthToken())

		//定时器
		botGroup.POST("/add_task", http.AddCronJob)
		botGroup.POST("/stop_task", http.StopTimeTask)
		botGroup.POST("/get_task_list", http.GetUserTaskList)
		botGroup.POST("/get_task_info", http.GetTaskInfo)
		//  @全体成员
		//botGroup.POST("/at_all_member", http.AtAllMember)

	}

	//websocket服务
	r.GET("/ws/bot", socket.BotWs)
	s.Router = r
}

func (s *Server) RunServer(address string) error {
	url := fmt.Sprintf("%s:%s", address, s.Port)
	return s.Router.Run(url)
}
