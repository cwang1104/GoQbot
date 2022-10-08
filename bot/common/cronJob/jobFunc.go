package cronJob

import (
	"log"
	"qbot/bot/common/tools"
)

func SendMsg(msgType string, sengToId int64, message string) error {
	msgSender := tools.NewMessageSender(msgType, sengToId, message)
	err := msgSender.SendMsg()
	if err != nil {
		log.Println("send msg failed", err)
		return err
	}
	return nil
}
