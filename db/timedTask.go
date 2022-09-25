package db

func AddTimedTask(task *TimedTaskModel) error {
	return db.Table("timed_task").Create(task).Error
}

func GetUserTaskList(userId, limit, offset int) (*[]TimedTaskModel, error) {
	var tasks []TimedTaskModel
	result := db.Table("timed_task").Where("created_id = ?", userId).Limit(limit).Offset(offset).Order("created_time desc").Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tasks, nil
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
