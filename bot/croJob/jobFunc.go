package croJob

import (
	"log"
	"qbot/bot/http"
)

func SendMsg(msgType string, sengToId int64, message string) error {
	msgSender := http.NewMessageSender(msgType, sengToId, message)
	err := msgSender.SendMsg()
	if err != nil {
		log.Println("send msg failed", err)
		return err
	}
	return nil
}
