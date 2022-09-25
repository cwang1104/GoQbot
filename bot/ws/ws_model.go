package ws

type MessageType struct {
	PostType    string `json:"post_type"`    // 上报类型：message request notice meta_event
	MessageType string `json:"message_type"` //消息类型：
	Time        int64  `json:"time"`
	SelfId      int64  `json:"self_id"`  //机器人qq号
	SubType     string `json:"sub_type"` //消息子类型  如果是好友则是 friend, 如果是群临时会话则是 group, 如果是在群中自身发送则是 group_self,
	// 正常群聊消息是 normal, 匿名消息是 anonymous, 系统提示 ( 如「管理员已禁止群内匿名聊天」 ) 是 notice
	Font       int64       `json:"font"`
	Sender     SenderModel `json:"sender"`
	MessageId  int64       `json:"message_id"`
	UserId     int64       `json:"user_id"`
	TargetId   int64       `json:"target_id"`
	Message    string      `json:"message"`
	RawMessage string      `json:"raw_message"`
	GroupId    int64       `json:"group_id"`
}

type SenderModel struct {
	Age      int64  `json:"age"`
	Area     string `json:"area"`
	Card     string `json:"card"`
	Level    string `json:"level"`
	NickName string `json:"nickname"`
	Role     string `json:"role"`
	Title    string `json:"title"`
	Sex      string `json:"sex"`
	UserId   int64  `json:"user_id"`
}

type SendWsMsgModel struct {
	Action string      `json:"action"`
	Params interface{} `json:"params"`
}
type SendGroupMsg struct {
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}
