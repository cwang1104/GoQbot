package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qbot/pkg/e"
	"qbot/pkg/logger"
	user2 "qbot/service/user"
)

func UserRegister(c *gin.Context) {
	var req user2.RegisterService
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Errorf("[params bindJson failed] [err:%v]", err)
		c.JSON(http.StatusBadRequest, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}
	res := req.AccountRegister()
	c.JSON(http.StatusOK, res)
}

func CheckLogin(c *gin.Context) {
	var req user2.RegisterService
	if err := c.BindJSON(&req); err != nil {
		logger.Log.Errorf("[params bindJson failed] [err:%v]", err)
		c.JSON(http.StatusBadRequest, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}
	res := req.CheckLogin()
	c.JSON(http.StatusOK, res)

}
