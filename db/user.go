package db

func AddUser(model *UserModel) error {
	return db.Model(&UserModel{}).Create(model).Error
}

func GetUserInfoByName(name string) (UserModel, error) {
	var user UserModel
	result := db.Model(UserModel{}).Where("user_name = ?", name).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
