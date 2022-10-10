package bot

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"qbot/bot/common/cronJob"
	"qbot/db"
	"qbot/middleware"
	"qbot/pkg/e"
	"qbot/pkg/utils"
	"time"
)

//TimedTaskService 定时任务service全参数
type TimedTaskService struct {
	Id             int          `json:"id"`
	CreatedId      int          `json:"created_id"`
	TaskName       string       `json:"task_name"`
	TimedStart     int          `json:"timed_start"`
	StartTime      int64        `json:"start_time"`
	TimedEnd       int          `json:"timed_end"`
	EndTime        int64        `json:"end_time"`
	TimingStrategy TimeStrategy `json:"timing_strategy"`
	TimerTypeId    int          `json:"timer_type_id"`
	SendType       string       `json:"send_type"`
	TaskExplain    string       `json:"task_explain"`
	SentContent    string       `json:"sent_content"`
	SendTo         int64        `json:"send_to"`
	Status         int          `json:"status"`
	CreatedTime    int64        `json:"created_time"`
	ModelPage
}

//TimeStrategy 时间策略
type TimeStrategy struct {
	Interval       int `json:"interval"` //分钟数 15 30 60 120
	TimeLimitStart int `json:"time_limit_start"`
	TimeLimitEnd   int `json:"time_limit_end"`
}

//ModelPage 获取列表页参数
type ModelPage struct {
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

//TimeTaskInfoResp 获取定时任务详情响应数据
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

//TaskListResp 获取定时任务列表响应数据
type TaskListResp struct {
	Id             int          `json:"id"`
	TaskName       string       `json:"task_name"`
	TimerTypeName  string       `json:"timer_type_name"`
	TimingStrategy TimeStrategy `json:"timing_strategy"`
	Status         int          `json:"status"`
	TaskExplain    string       `json:"task_explain"`
	CreatedTime    int64        `json:"created_time"`
}

//AddCronJob 添加定时任务
func (t *TimedTaskService) AddCronJob(ctx *gin.Context) (res gin.H) {
	payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)
	timingStrategyJson, err := json.Marshal(t.TimingStrategy)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_UNMARSHAL_JSON, err)
		return
	}

	taskInfo := db.TimedTaskModel{
		CreatedId:      payload.UID,
		TaskName:       t.TaskName,
		TimedStart:     t.TimedStart,
		StartTime:      t.StartTime,
		TimedEnd:       t.TimedEnd,
		EndTime:        t.EndTime,
		TimingStrategy: string(timingStrategyJson),
		TimerTypeId:    t.TimerTypeId,
		SendType:       t.SendType,
		TaskExplain:    t.TaskExplain,
		SentContent:    t.SentContent,
		SendTo:         t.SendTo,
		Status:         1,
		CreatedTime:    time.Now().Unix(),
	}

	err = db.AddTimedTask(&taskInfo)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_DATABASE_CREATE_FAIL, err)
		return
	}

	newTaskInfo, err := db.GetTaskInfoByNameAndUserId(taskInfo.TaskName, taskInfo.CreatedId)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_DATABASE_QUERY_FAIL, err)
		return
	}

	timeJob, err := cronJob.NewCronJob(newTaskInfo)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_CREATE_CRONJOB_FAIL, err)
		return
	}
	timeJob.StartCronJob()
	res = e.SuccessResponse()
	return
}

//StopTimeTask 停止定时任务
func (t *TimedTaskService) StopTimeTask(ctx *gin.Context) (res gin.H) {
	payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)

	err := db.UpdateTaskStatusByUser(3, t.Id, payload.UID)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_DATABASE_UPDATE_FAIL, err)
		return
	}

	taskName := cronJob.GetTimeTaskName(t.TaskName, t.Id)
	timeTask, exist := cronJob.GetTimedTask(taskName)

	if !exist {
		res = e.ErrorResponse(e.ERROR_CROBJOB_NOT_EXIST, nil)
		return
	}
	timeTask.StopCronJob()

	res = e.SuccessResponse()
	return
}

//GetUserTaskList 获取当前用户的定时任务列表
func (t *TimedTaskService) GetUserTaskList(ctx *gin.Context) (res gin.H) {
	payload := ctx.MustGet(middleware.AuthorizationPayloadKey).(*utils.Claims)
	offset := (t.CurrentPage - 1) * t.PageSize

	//根据前端传入的status查询对应状态的定时任务
	var tasks *[]db.GetTaskInfoModel
	var err error
	if t.Status == 4 {
		tasks, err = db.GetUserTaskList(payload.UID, t.PageSize, offset)
	} else {
		fmt.Println(t.Status)
		tasks, err = db.GetUserTaskListStatus(payload.UID, t.PageSize, offset, t.Status)
	}
	if err != nil && err != io.EOF {
		res = e.ErrorResponse(e.ERROR_DATABASE_QUERY_FAIL, err)
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

	data := gin.H{
		"tasks": tasksResp,
	}
	res = e.SuccessResponseWithData(data)
	return
}

//GetTaskInfo 根据id获取定时器详情
func (t *TimedTaskService) GetTaskInfo() (res gin.H) {

	task, err := db.GetTaskInfoById(t.Id)
	if err != nil && err != io.EOF {
		e.ErrorResponse(e.ERROR_DATABASE_QUERY_FAIL, err)
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

	data := gin.H{
		"taskInfo": taskResp,
	}
	res = e.SuccessResponseWithData(data)
	return
}
