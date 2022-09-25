package croJob

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"log"
	"qbot/db"
	"qbot/utils"
	"sync"
	"time"
)

var (
	TimedTaskList = map[string]*CronJob{}
	TimeLock      sync.Mutex
)

type CronJob struct {
	cro            *cron.Cron //定时器
	TaskId         int
	CreatedId      int
	TaskName       string
	TimedStart     int
	StartTime      int64
	TimedEnd       int
	EndTime        int64
	Content        string
	SendTo         int
	TimingStrategy *TimeStrategy
	TimerType      int
	Status         int
}

//TimeStrategy 时间策略
type TimeStrategy struct {
	Interval       int `json:"interval"` //分钟数 15 30 60 120
	TimeLimitStart int `json:"time_limit_start"`
	TimeLimitEnd   int `json:"time_limit_end"`
}

func NewCronJob(timedTask *db.TimedTaskModel) (*CronJob, error) {

	//将数据库存的TimingStrategy json字符串转成结构体指针存入
	var timeStrategy TimeStrategy
	err := json.Unmarshal([]byte(timedTask.TimingStrategy), &timeStrategy)
	if err != nil {
		log.Println("TimingStrategy json string", timedTask.TimingStrategy)
		log.Println("unmarshal json string failed", err)
		return nil, err
	}

	return &CronJob{
		cro:            cron.New(),
		TaskId:         timedTask.Id,
		CreatedId:      timedTask.CreatedId,
		TaskName:       timedTask.TaskName,
		TimedStart:     timedTask.TimedStart,
		StartTime:      timedTask.StartTime,
		TimedEnd:       timedTask.TimedEnd,
		EndTime:        timedTask.EndTime,
		Content:        timedTask.SentContent,
		SendTo:         timedTask.SendTo,
		TimingStrategy: &timeStrategy,
		TimerType:      timedTask.TimerType,
		Status:         timedTask.Status,
	}, nil
}

//StartCronJob 开启定时器
func (c *CronJob) StartCronJob() {

	spec := utils.GetInternalSpec(c.TimingStrategy.Interval, c.TimingStrategy.TimeLimitStart, c.TimingStrategy.TimeLimitEnd)
	log.Println("--------cron func", spec)

	_ = c.cro.AddFunc(spec, func() {
		fmt.Println("task start")
	})

	//添加进TimeTaskList
	taskName := utils.GetTimeTaskName(c.TaskName, c.TaskId)
	err := AddTimedTask(c)
	if err != nil {
		log.Println("AddTimedTask failed", err)
		return
	}
	//检验是否 定时启动和定时关闭
	if c.TimedStart == 1 {
		err := db.UpdateTaskStatus(1, c.TaskId)
		if err != nil {
			log.Println("UpdateTaskStatus failed", err)
			return
		}
		for {
			if time.Now().Unix() == c.StartTime {
				c.cro.Start()
				err := db.UpdateTaskStatus(2, c.TaskId)
				if err != nil {
					log.Println("UpdateTaskStatus failed", err)
					return
				}
			}
		}
	} else {
		c.cro.Start()
		err := db.UpdateTaskStatus(2, c.TaskId)
		if err != nil {
			log.Println("UpdateTaskStatus failed", err)
			return
		}
	}

	//定时关闭
	if c.TimedEnd == 1 {
		for {
			if time.Now().Unix() == c.EndTime {
				c.cro.Stop()
				DelTimedTask(taskName)
				err := db.UpdateTaskStatus(3, c.TaskId)
				if err != nil {
					log.Println("UpdateTaskStatus failed", err)
					return
				}
			}
		}
	}
}

func (c *CronJob) StopCronJob() {
	taskName := utils.GetTimeTaskName(c.TaskName, c.TaskId)
	DelTimedTask(taskName)
	c.cro.Stop()
	err := db.UpdateTaskStatus(3, c.TaskId)
	if err != nil {
		log.Println("db UpdateTaskStatus failed", err)
		return
	}
}

// AddTimedTask 添加定时任务到map
func AddTimedTask(cronJob *CronJob) error {
	TimeLock.Lock()
	taskName := utils.GetTimeTaskName(cronJob.TaskName, cronJob.TaskId)
	_, exist := GetTimedTask(taskName)
	if exist {
		TimeLock.Unlock()
		return errors.New("task is exist")
	}
	TimedTaskList[taskName] = cronJob
	TimeLock.Unlock()
	return nil
}

//GetTimedTask 获取定时任务
func GetTimedTask(taskName string) (*CronJob, bool) {
	TimeLock.Lock()
	cronJob, exist := TimedTaskList[taskName]
	TimeLock.Unlock()
	return cronJob, exist
}

//DelTimedTask 删除定时任务
func DelTimedTask(taskName string) {
	TimeLock.Lock()
	delete(TimedTaskList, taskName)
	TimeLock.Unlock()
}
