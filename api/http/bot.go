package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	http2 "qbot/bot/http"
)

type AtAllReq struct {
	GroupId int64 `json:"group_id"`
}

func AtAllMember(c *gin.Context) {

	deal := http2.NewMemberDeal(527388259, 0, false)
	err := deal.AtAllMember()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "al all member failed" + err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "成功",
	})
}
