package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"qbot/db"
	utils2 "qbot/pkg/utils"
	"time"
)

type UserRegisterReq struct {
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func UserRegister(c *gin.Context) {
	var req UserRegisterReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}

	hashPassword, err := utils2.HashPassword(req.Password)
	if err != nil {
		log.Println("HashPassword failed", err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "HashPassword failed:" + err.Error(),
		})
		return
	}

	user := db.UserModel{
		UserName:    req.UserName,
		Password:    hashPassword,
		CreatedTime: time.Now().Unix(),
	}

	err = db.AddUser(&user)
	if err != nil {
		log.Println("db add user failed", err)
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "db add user failed:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func CheckLogin(c *gin.Context) {
	var req UserRegisterReq
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "req failed " + err.Error(),
		})
		return
	}

	userInfo, err := db.GetUserInfoByName(req.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "db failed" + err.Error(),
		})
		return
	}

	err = utils2.CheckPassword(req.Password, userInfo.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2002,
			"msg":  "password error",
		})
		return
	}

	token, err := utils2.NewJWT().CreateToken(userInfo.Id, userInfo.UserName)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2002,
			"msg":  "create token failed" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"uid":      userInfo.Id,
			"userName": userInfo.UserName,
			"token":    token,
		},
	})
}
