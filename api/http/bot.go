package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"qbot/api/middleware"
	"qbot/bot/croJob"
	"qbot/db"
	"qbot/utils"
	"time"
)

type TimeTaskInfoResp struct {
	Id             int          `json:"id"`
	CreatedId      int          `json:"created_id"`
	TaskName       string       `json:"task_name"`
	TimedStart     int          `json:"timed_start"`
	StartTime      int64        `json:"start_time"`
	TimedEnd       int          `json:"timed_end"`
	EndTime        int64        `json:"end_time"`
	TimingStrategy TimeStrategy `json:"timing_strategy"`
	TimerTypeId    int          `json:"timer_type_id"`
	TimerTypeName  string       `json:"timer_type_name"`
	TaskExplain    string       `json:"task_explain"`
	SendType       string       `json:"send_type"`
	SentContent    string       `json:"sent_content"`
	SendTo         int64        `json:"send_to"`
	Status         int          `json:"status"`
	CreatedTime    int64        `json:"created_time"`
}

type TaskListResp struct {
	Id             int          `json:"id"`
	TaskName       string       `json:"task_name"`
	TimerTypeName  string       `json:"timer_type_name"`
	TimingStrategy TimeStrategy `json:"timing_strategy"`
	Status         int          `json:"status"`
	TaskExplain    string       `json:"task_explain"`
	CreatedTime    int64        `json:"created_time"`
}

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
	TimerTypeId int    `json:"timer_type_id"`
	SendType    string `json:"send_type"`
	TaskExplain string `json:"task_explain"`
	SentContent string `json:"sent_content"`
	SendTo      int64  `json:"send_to"`
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
		TimerTypeId:    req.TimerTypeId,
		SendType:       req.SendType,
		TaskExplain:    req.TaskExplain,
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

	//taskName := utils.GetTimeTaskName(req.TaskName, req.TaskId)
	//cronJob, exist := croJob.GetTimedTask(taskName)
	//if !exist {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"code": 500,
	//		"msg":  "task is not exist timeTask" + err.Error(),
	//	})
	//	return
	//}
	//cronJob.StopCronJob()

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

type GetReq struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	Status      int `json:"status"`
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
	log.Println("---------------")
	fmt.Printf("%+v\n", req)
	payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)
	offset := (req.CurrentPage - 1) * req.PageSize
	var tasks *[]db.GetTaskInfoModel
	var err error
	if req.Status == 4 {
		tasks, err = db.GetUserTaskList(payload.UID, req.PageSize, offset)
	} else {
		fmt.Println(req.Status)
		tasks, err = db.GetUserTaskListStatus(payload.UID, req.PageSize, offset, req.Status)
	}
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 3001,
			"msg":  "GetUserTaskList failed : " + err.Error(),
		})
		return
	}
	var tasksResp []TaskListResp
	for _, task := range *tasks {
		var strategy TimeStrategy
		_ = json.Unmarshal([]byte(task.TimingStrategy), &strategy)
		taskResp := TaskListResp{
			Id:             task.Id,
			TaskName:       task.TaskName,
			TimingStrategy: strategy,
			TaskExplain:    task.TaskExplain,
			Status:         task.Status,
			TimerTypeName:  task.TimerTypeName,
			CreatedTime:    task.CreatedTime,
		}
		tasksResp = append(tasksResp, taskResp)
	}

	fmt.Printf("%+v\n", tasksResp)
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"tasks": tasksResp,
		},
	})
}

type GetTaskInfoReq struct {
	TaskId int `json:"task_id" binding:"required,min=1"`
}

func GetTaskInfo(c *gin.Context) {
	var req GetTaskInfoReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}
	//payload := c.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)

	task, err := db.GetTaskInfoById(req.TaskId)
	if err != nil && err != io.EOF {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "db failed" + err.Error(),
		})
		return
	}

	var strategy TimeStrategy
	_ = json.Unmarshal([]byte(task.TimingStrategy), &strategy)
	taskResp := TimeTaskInfoResp{
		Id:             task.Id,
		CreatedId:      task.CreatedId,
		TaskName:       task.TaskName,
		TimedStart:     task.TimedStart,
		StartTime:      task.StartTime,
		TimedEnd:       task.TimedEnd,
		EndTime:        task.EndTime,
		TimingStrategy: strategy,
		TimerTypeId:    task.TimerTypeId,
		TimerTypeName:  task.TimerTypeName,
		SendType:       task.SendType,
		TaskExplain:    task.TaskExplain,
		SentContent:    task.SentContent,
		SendTo:         task.SendTo,
		Status:         task.Status,
		CreatedTime:    task.CreatedTime,
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"taskInfo": taskResp,
		},
	})

}
