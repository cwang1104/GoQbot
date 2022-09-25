package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"qbot/api/middleware"
	"qbot/bot/croJob"
	"qbot/db"
	"qbot/utils"
	"time"
)

type AddTimedTaskReq struct {
	TaskName       string `json:"task_name"`
	TimedStart     int    `json:"timed_start"`
	StartTime      int64  `json:"start_time"`
	TimedEnd       int    `json:"timed_end"`
	EndTime        int64  `json:"end_time"`
	TimingStrategy struct {
		Interval       int `json:"interval"` //分钟数 15 30 60 120
		TimeLimitStart int `json:"time_limit_start"`
		TimeLimitEnd   int `json:"time_limit_end"`
	} `json:"timing_strategy"`
	TimerType   int    `json:"timer_type"`
	SentContent string `json:"sent_content"`
	SendTo      int    `json:"send_to"`
}

//TimeStrategy 时间策略
type TimeStrategy struct {
	Interval       int `json:"interval"` //分钟数 15 30 60 120
	TimeLimitStart int `json:"time_limit_start"`
	TimeLimitEnd   int `json:"time_limit_end"`
}

func AddCronJob(c *gin.Context) {
	var req AddTimedTaskReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}

	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)

	timingStrategyJson, err := json.Marshal(req.TimingStrategy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 2001,
			"msg":  "TimingStrategy json failed:" + err.Error(),
		})
	}

	taskInfo := db.TimedTaskModel{
		CreatedId:      payload.UID,
		TaskName:       req.TaskName,
		TimedStart:     req.TimedStart,
		StartTime:      req.StartTime,
		TimedEnd:       req.TimedEnd,
		EndTime:        req.EndTime,
		TimingStrategy: string(timingStrategyJson),
		TimerType:      req.TimerType,
		SentContent:    req.SentContent,
		SendTo:         req.SendTo,
		Status:         1,
		CreatedTime:    time.Now().Unix(),
	}

	err = db.AddTimedTask(&taskInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "db add time task failed" + err.Error(),
		})
		return
	}

	newTaskInfo, err := db.GetTaskInfoByNameAndUserId(taskInfo.TaskName, taskInfo.CreatedId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "db get time task failed" + err.Error(),
		})
		return
	}

	cronJob, err := croJob.NewCronJob(newTaskInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "NewCronJob failed" + err.Error(),
		})
		return
	}
	cronJob.StartCronJob()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

type StopTimeTaskReq struct {
	TaskName string `json:"task_name"`
	TaskId   int    `json:"task_id"`
}

func StopTimeTask(c *gin.Context) {
	var req StopTimeTaskReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}

	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)

	err := db.UpdateTaskStatusByUser(3, req.TaskId, payload.UID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "db UpdateTaskStatusByUser failed" + err.Error(),
		})
		return
	}

	taskName := utils.GetTimeTaskName(req.TaskName, req.TaskId)
	cronJob, exist := croJob.GetTimedTask(taskName)
	if !exist {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "task is not exist" + err.Error(),
		})
		return
	}
	cronJob.StopCronJob()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

type GetReq struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

func GetUserTaskList(c *gin.Context) {
	var req GetReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)
	offset := (req.CurrentPage - 1) * req.PageSize
	tasks, err := db.GetUserTaskList(payload.UID, req.PageSize, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 3001,
			"msg":  "GetUserTaskList failed : " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"tasks": tasks,
		},
	})
}
