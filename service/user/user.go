package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"qbot/db"
	"qbot/pkg/e"
	utils2 "qbot/pkg/utils"
	"time"
)

type RegisterService struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

//AccountRegister 注册
func (r *RegisterService) AccountRegister() gin.H {
	hashPassword, err := utils2.HashPassword(r.Password)
	if err != nil {
		log.Println("HashPassword failed", err.Error())
		res := e.ErrorResponse(e.ERROR_FAIL_ENCRYPTION, err)
		return res
	}

	user := db.UserModel{
		UserName:    r.UserName,
		Password:    hashPassword,
		CreatedTime: time.Now().Unix(),
	}

	err = db.AddUser(&user)
	if err != nil {
		log.Println("db add user failed", err)
		res := e.ErrorResponse(e.ERROR_DATABASE, err)
		return res
	}
	return e.SuccessResponse()
}

//CheckLogin 登录
func (r *RegisterService) CheckLogin() (res gin.H) {

	//数据库中查询密码
	userInfo, err := db.GetUserInfoByName(r.UserName)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_DATABASE, err)
		return
	}

	//查验密码是否出错
	err = utils2.CheckPassword(r.Password, userInfo.Password)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_CHECK_PASSWORD_FAIL, err)
		return
	}

	//创建token
	//todo: 将token存入缓存防止多端登录
	token, err := utils2.NewJWT().CreateToken(userInfo.Id, userInfo.UserName)
	if err != nil {
		res = e.ErrorResponse(e.ERROR_GENERATE_TOKEN, err)
		return
	}

	data := gin.H{
		"uid":      userInfo.Id,
		"userName": userInfo.UserName,
		"token":    token,
	}
	res = e.SuccessResponseWithData(data)
	return
}
