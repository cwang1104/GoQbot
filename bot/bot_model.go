package bot

// 成员信息
type MemberInfo struct {
	GroupId         int64  `json:"group_id"`
	UserId          int64  `json:"user_id"`
	NickName        string `json:"nick_name"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int32  `json:"age"`
	Area            string `json:"area"`
	JoinTime        int32  `json:"join_time"`
	LastSentTime    int32  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int64  `json:"title_expire_time"`
	CardChangeAble  bool   `json:"card_change_able"`
	ShutUpTimestamp int64  `json:"shut_up_timestamp"`
}

//发送消息api请求参数
type SendPrivatePMsg struct {
	UserID     int64  `json:"user_id"`
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

//获取群成员api请求参数
type GetMemberList struct {
	GroupId int64 `json:"group_id"`
	NoCache bool  `json:"no_cache"`
}

type SendGroupMsg struct {
	GroupID    int64  `json:"group_id"`
	Message    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}
