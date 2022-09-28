package croJob

import (
	"encoding/json"
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
	SentContent    string
	SendTo         int64
	TimingStrategy *TimeStrategy
	TimerTypeId    int
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
		SentContent:    timedTask.SentContent,
		SendTo:         timedTask.SendTo,
		TimingStrategy: &timeStrategy,
		TimerTypeId:    timedTask.TimerTypeId,
		Status:         timedTask.Status,
	}, nil
}

//StartCronJob 开启定时器
func (c *CronJob) StartCronJob() {

	spec, err := utils.GetInternalSpec(c.TimingStrategy.Interval, c.TimingStrategy.TimeLimitStart, c.TimingStrategy.TimeLimitEnd)
	if err != nil {
		log.Println("GetInternalSpec failed", err)
		return
	}
	log.Println("--------cron func", spec)

	err = c.cro.AddFunc(spec, func() {
		SendMsg("private", c.SendTo, c.SentContent)
	})
	if err != nil {
		log.Println("AddFunc", err)
		return
	}
	//添加进TimeTaskList
	taskName := utils.GetTimeTaskName(c.TaskName, c.TaskId)
	AddTimedTask(c)

	//检验是否 定时启动和定时关闭
	if c.TimedStart == 1 {
		go func() {
			for {
				if time.Now().Unix() == c.StartTime {
					c.cro.Start()
					err := db.UpdateTaskStatus(2, c.TaskId)
					if err != nil {
						log.Println("UpdateTaskStatus failed", err)
						return
					}
					log.Println("任务开始", c.TaskName)
					break
				}
			}
		}()

	} else {
		c.cro.Start()
		err := db.UpdateTaskStatus(2, c.TaskId)
		if err != nil {
			log.Println("UpdateTaskStatus failed", err)
			return
		}
		log.Println("任务开始", c.TaskName)
	}

	//定时关闭
	if c.TimedEnd == 1 {
		go func() {
			for {
				if time.Now().Unix() == c.EndTime {
					c.cro.Stop()
					DelTimedTask(taskName)
					err := db.UpdateTaskStatus(3, c.TaskId)
					if err != nil {
						log.Println("UpdateTaskStatus failed", err)
						return
					}
					log.Println("任务结束", c.TaskName)
					break
				}
			}
		}()
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
func AddTimedTask(cronJob *CronJob) {
	TimeLock.Lock()
	taskName := utils.GetTimeTaskName(cronJob.TaskName, cronJob.TaskId)
	TimedTaskList[taskName] = cronJob
	TimeLock.Unlock()
	return
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
