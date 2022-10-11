package main

import (
	"log"
	"qbot/api"
	"qbot/bot/common/cronJob"
	"qbot/db"
	"qbot/pkg/logger"
	"qbot/pkg/utils"
)

func init() {
	err := db.DbInit()
	if err != nil {
		logger.Log.Errorf("db init failed,err: %v", err)
		panic("db init failed")
	}

	logger.Log.Infof("db init success")

	err = cronJob.TimeTaskInit()
	if err != nil {
		logger.Log.Errorf("bot time task init failed,err: %v", err)
		panic("time task init failed")
	}
	logger.Log.Infof("bot init success")
}

func main() {
	server := api.NewServer(utils.GlobalConf.Server.Port)
	err := server.RunServer(utils.GlobalConf.Server.Ip)
	if err != nil {
		log.Println("run server failed", err)
		return
	}
	logger.Log.Infof("run api server success")
}
