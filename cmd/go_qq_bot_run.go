package main

import (
	"log"
	"qbot/api"
	"qbot/db"
	"qbot/utils"
)

func init() {
	err := db.DbInit()
	if err != nil {
		log.Println("DB init failed", err)
		return
	}
	log.Println("db init success")
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
