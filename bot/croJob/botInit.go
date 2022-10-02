package croJob

import (
	"log"
	"qbot/db"
)

func TimeTaskInit() error {
	tasks, err := db.GetRunningTask()
	if err != nil {
		log.Println("init bot get running task failed", err)
		return err
	}

	if len(*tasks) == 0 {
		log.Println("当前无运行中任务")
		return nil
	}

	for _, task := range *tasks {
		cronJob, err := NewCronJob(&task)
		if err != nil {
			log.Println("new cronJob failed", err)
			return err
		}
		cronJob.StartCronJob()
	}
	return nil
}
