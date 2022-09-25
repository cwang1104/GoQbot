package db

func AddUser(model *UserModel) error {
	return db.Table("user").Create(model).Error
}

func GetUserInfoByName(name string) (UserModel, error) {
	var user UserModel
	result := db.Table("user").Where("user_name = ?", name).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
