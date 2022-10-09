package http

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"qbot/pkg/e"
	user2 "qbot/service/user"
)

func UserRegister(c *gin.Context) {
	var req user2.RegisterService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusBadRequest, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}
	res := req.AccountRegister()
	c.JSON(http.StatusOK, res)
}

func CheckLogin(c *gin.Context) {
	var req user2.RegisterService
	if err := c.BindJSON(&req); err != nil {
		log.Println("bindJson failed", err)
		c.JSON(http.StatusBadRequest, e.ErrorResponse(e.INVALID_PARAMS, err))
		return
	}
	res := req.CheckLogin()
	c.JSON(http.StatusOK, res)

}
