package db

type UserModel struct {
	Id          int    `json:"id"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
	CreatedTime int64  `json:"created_time"`
}

type TimedTaskModel struct {
	Id             int    `json:"id"`
	CreatedId      int    `json:"created_id"`
	TaskName       string `json:"task_name"`
	TimedStart     int    `json:"timed_start"`
	StartTime      int64  `json:"start_time"`
	TimedEnd       int    `json:"timed_end"`
	EndTime        int64  `json:"end_time"`
	TimingStrategy string `json:"timing_strategy"`
	TimerTypeId    int    `json:"timer_type_id"`
	SendType       string `json:"send_type"`
	TaskExplain    string `json:"task_explain"`
	SentContent    string `json:"sent_content"`
	SendTo         int64  `json:"send_to"`
	Status         int    `json:"status"`
	CreatedTime    int64  `json:"created_time"`
}

type TimerTypeModel struct {
	Id          int    `json:"id"`
	TypeName    string `json:"type_name"`
	TypeExplain string `json:"type_explain"`
	CreatedTime int64  `json:"created_time"`
}

type GetTaskInfoModel struct {
	Id             int    `json:"id"`
	CreatedId      int    `json:"created_id"`
	TaskName       string `json:"task_name"`
	TimedStart     int    `json:"timed_start"`
	StartTime      int64  `json:"start_time"`
	TimedEnd       int    `json:"timed_end"`
	EndTime        int64  `json:"end_time"`
	TimingStrategy string `json:"timing_strategy"`
	TimerTypeId    int    `json:"timer_type_id"`
	SendType       string `json:"send_type"`
	TaskExplain    string `json:"task_explain"`
	SentContent    string `json:"sent_content"`
	SendTo         int64  `json:"send_to"`
	Status         int    `json:"status"`
	CreatedTime    int64  `json:"created_time"`
	TimerTypeName  string `json:"timer_type_name"`
}
