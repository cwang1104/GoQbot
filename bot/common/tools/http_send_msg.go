package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"qbot/pkg/logger"
)

type MessageSender struct {
	MessageType string `json:"message_type"`
	UserId      int64  `json:"user_id"`
	GroupId     int64  `json:"group_id"`
	Message     string `json:"message"`
	AutoEscape  bool   `json:"auto_escape"`
}

func NewMessageSender(messageType string, sendToId int64, message string) *MessageSender {

	var msgSender MessageSender
	msgSender.MessageType = messageType
	msgSender.Message = message
	msgSender.AutoEscape = false
	if messageType == "private" {
		msgSender.UserId = sendToId
	} else if messageType == "group" {
		msgSender.GroupId = sendToId
	}
	return &msgSender
}

func (m *MessageSender) SendMsg() error {
	b, err := json.Marshal(m)
	if err != nil {
		logger.Log.Errorf("[json marshal failed][err:%v][data:%+v]", err, m)
		return err
	}
	url := fmt.Sprintf("%s/send_msg", CqHttpBaseUrl)
	content, err := HttpPost(url, b)
	if err != nil {
		logger.Log.Errorf("[http post failed][err:%v]", err)
		return err
	}
	log.Println(string(content))
	return nil
}
