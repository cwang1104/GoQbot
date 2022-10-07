package main

import (
	"log"
	"qbot/api"
	"qbot/bot/croJob"
	"qbot/db"
	"qbot/utils"
)

func init() {
	err := db.DbInit()
	if err != nil {
		log.Println("DB init failed", err)
		panic("db init failed")
	}
	log.Println("db init success")

	err = croJob.TimeTaskInit()
	if err != nil {
		log.Println("bot init failed", err)
		panic("time task init failed")
	}

}

func main() {
	server := api.NewServer(utils.GlobalConf.Server.Port)
	err := server.RunServer(utils.GlobalConf.Server.Ip)
	if err != nil {
		log.Println("run server failed", err)
		return
	}
	log.Println("run server success...")
}
