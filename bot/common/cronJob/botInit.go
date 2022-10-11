package cronJob

import (
	"qbot/db"
	"qbot/pkg/logger"
)

func TimeTaskInit() error {
	logger.Log.Infof("[start init bot time task...]")
	tasks, err := db.GetRunningTask()
	if err != nil {
		logger.Log.Errorf("[init bot get running task failed][err:%v]", err)
		return err
	}

	if len(*tasks) == 0 {
		logger.Log.Infof("[定时任务初始化完成,当前无定时任务...]")
		return nil
	}

	for _, task := range *tasks {
		cronJob, err := NewCronJob(&task)
		if err != nil {
			logger.Log.Errorf("[startrunning task failed][err:%v]", err)
			return err
		}
		cronJob.StartCronJob()
	}
	return nil
}
