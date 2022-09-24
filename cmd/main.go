package main

import (
	"log"
	"qbot/api"
	"qbot/db"
)

func init() {
	err := db.DbInit()
	if err != nil {
		log.Println("DBinit failed", err)
		return
	}
	log.Println("db init success")
}

func main() {
	server := api.NewServer("12000")
	err := server.RunServer("127.0.0.1")
	if err != nil {
		log.Println("run server failed", err)
		return
	}
	log.Println("run server success...")
}
