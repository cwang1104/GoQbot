package member_change_notice

import (
	"fmt"
	"qbot/bot/common/tools"
	"qbot/pkg/logger"
)

func MemberJoinNotice(message *tools.MessageType) {
	//群成员增加
	if message.NoticeType == "group_increase" {
		msg := fmt.Sprintf("欢迎[CQ:at,qq=%d]进群，请遵守群规范，关注群公告~\n来自bot", message.UserId)
		sendMsg := tools.GetGroupMsgStruct(msg, message.GroupId)
		tools.MsgChan <- sendMsg
		logger.Log.Infof("[%d用户加入了群组%d]", message.UserId, message.GroupId)
	}
}
