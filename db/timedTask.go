package db

import "fmt"

func AddTimedTask(task *TimedTaskModel) error {
	return db.Table("timed_task").Create(task).Error
}

func GetUserTaskList(userId, limit, offset int) (*[]GetTaskInfoModel, error) {
	var resp []GetTaskInfoModel
	sql := fmt.Sprintf("SELECT timed_task.*,timer_type.`type_name` AS timer_type_name \nFROM timed_task,timer_type \nWHERE timed_task.`created_id` = %d\nAND timed_task.`timer_type_id` = timer_type.`id`\nORDER BY timed_task.`created_time` DESC,timed_task.`id`\nLIMIT %d\nOFFSET %d;",
		userId, limit, offset)
	err := db.Raw(sql).Scan(&resp).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v\n", resp)
	return &resp, nil
}

func GetUserTaskListStatus(userId, limit, offset, status int) (*[]GetTaskInfoModel, error) {
	var resp []GetTaskInfoModel
	sql := fmt.Sprintf("SELECT timed_task.*,timer_type.`type_name` AS timer_type_name \nFROM timed_task,timer_type \nWHERE timed_task.`created_id` = %d\nAND timer_type.`id` = timed_task.`timer_type_id` \nAND timed_task.`status` = %d\nORDER BY timed_task.`created_time` DESC,timed_task.`id`\nLIMIT %d\nOFFSET %d;",
		userId, status, limit, offset)
	fmt.Println(sql)
	err := db.Raw(sql).Scan(&resp).Error
	if err != nil {
		return nil, err
	}
	fmt.Printf("%+v", resp)
	return &resp, nil
}

func GetRunningTask() (*[]TimedTaskModel, error) {
	var req []TimedTaskModel
	result := db.Table("timed_task").Where("status = 2").Find(&req)
	if result.Error != nil {
		return nil, result.Error
	}
	return &req, nil
}

func GetAllTaskList() (*[]TimedTaskModel, error) {
	var tasks []TimedTaskModel
	result := db.Table("timed_task").Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tasks, nil
}

func UpdateTaskStatus(status, id int) error {
	return db.Exec("UPDATE timed_task SET `status` = ? WHERE id = ?", status, id).Error
}

func UpdateTaskStatusByUser(status, taskId, createdId int) error {
	return db.Exec("UPDATE timed_task SET `status` = ? WHERE id = ? AND created_id = ?", status, taskId, createdId).Error
}

func GetTaskInfoByNameAndUserId(taskName string, createdId int) (*TimedTaskModel, error) {
	var task TimedTaskModel
	result := db.Table("timed_task").Where("created_id = ?", createdId).Where("task_name = ?", taskName).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

func GetTaskInfoById(taskId int) (*GetTaskInfoModel, error) {
	var resp GetTaskInfoModel
	sql := fmt.Sprintf("SELECT timed_task.*,timer_type.`type_name` \nFROM timed_task,timer_type \nWHERE timed_task.`id` = %d\nAND timed_task.`timer_type_id` = timer_type.`id`;", taskId)
	err := db.Raw(sql).Scan(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
