package ws

//SendWsMsgModel 发送消息 消息体结构
type SendWsMsgModel struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}

//SendGroupMsg 发送群消息消息 params结构
type SendGroupMsg struct {
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

//SendPrivateMsg 发送私聊消息消息 params结构
type SendPrivateMsg struct {
	UserId     int64  `json:"user_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

//GetGroupMsgStruct 获取发送qq群消息api结构
func GetGroupMsgStruct(message string, groupId int64) *SendWsMsgModel {
	msg := SendGroupMsg{
		GroupID:    groupId,
		Message:    message,
		AutoEscape: false,
	}

	wsMsg := SendWsMsgModel{
		Action: "send_group_msg",
		Params: msg,
	}
	return &wsMsg
}

//GetPrivateMsgStruct 获取发送私聊消息api结构
func GetPrivateMsgStruct(message string, userId int64) *SendWsMsgModel {
	msg := SendPrivateMsg{
		UserId:     userId,
		Message:    message,
		AutoEscape: false,
	}

	wsMsg := SendWsMsgModel{
		Action: "send_private_msg",
		Params: msg,
	}
	return &wsMsg
}
