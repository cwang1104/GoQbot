package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"qbot/bot/common/cronJob"
	"qbot/pkg/e"
	"qbot/service/bot"
)

//AddCronJob 添加定时任务
func AddCronJob(c *gin.Context) {
	var req bot.TimedTaskService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}

	res := req.AddCronJob(c)
	c.JSON(http.StatusOK, res)

}

//StopTimeTask 手动停止定时任务
/*
	参数：
		task_name string
		id int
*/
func StopTimeTask(c *gin.Context) {
	var req bot.TimedTaskService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}

	res := req.StopTimeTask(c)
	c.JSON(http.StatusOK, res)
}

//GetUserTaskList 获取用户定时任务列表
func GetUserTaskList(c *gin.Context) {
	var req bot.TimedTaskService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}

	res := req.GetUserTaskList(c)
	c.JSON(http.StatusOK, res)

}

//GetTaskInfo 获取定时任务详情
/*
	参数：
		current_page int
		page_size int
		id int
*/
func GetTaskInfo(c *gin.Context) {
	var req bot.TimedTaskService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}
	//payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)
	res := req.GetTaskInfo()
	c.JSON(http.StatusOK, res)

}

func ShowCronMap(c *gin.Context) {
	for k, v := range cronJob.TimedTaskList {
		fmt.Printf("key = %v, value = %v\n", k, v)
	}
}
