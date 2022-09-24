package db

func AddTimedTask(task *TimedTaskModel) error {
	return db.Table("timed_task").Create(task).Error
}

func GetUserTaskList(userId int) (*[]TimedTaskModel, error) {
	var tasks []TimedTaskModel
	result := db.Table("timed_task").Where("created_id = ?", userId).Find(&tasks)
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
