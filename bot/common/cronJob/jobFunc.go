package cronJob

import (
	"qbot/bot/common/tools"
	"qbot/pkg/logger"
)

func SendMsg(msgType string, sengToId int64, message string) error {
	msgSender := tools.NewMessageSender(msgType, sengToId, message)
	err := msgSender.SendMsg()
	if err != nil {
		logger.Log.Errorf("[sendMsg failed][err:%v]", err)
		return err
	}
	return nil
}
