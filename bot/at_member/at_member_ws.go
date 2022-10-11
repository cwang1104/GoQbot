package at_member

import (
	"qbot/bot/common/tools"
	"qbot/pkg/logger"
)

func AtMeFunc(message *tools.MessageType) {
	msg := "功能菜单：\n 1-查询天气 (输入 /天气)\n 2-艾特全体 (输入 /艾特全体)\n"
	sendInfo := tools.GetGroupMsgStruct(msg, message.GroupId)
	tools.MsgChan <- sendInfo
	logger.Log.Infof("[user %s use AtMeFunction failed]", message.Sender.NickName)
}

// AtAllMember @全体成员
func AtAllMember(message *tools.MessageType) {
	if message.Sender.Role == "member" {
		sendMsg := tools.GetGroupMsgStruct("此功能仅群主及管理员可用", message.GroupId)
		tools.MsgChan <- sendMsg
		logger.Log.Infof("[user %s use AtAllMember]", message.Sender.NickName)
		return
	}

	//获取群成员的信息，构建qq号切片
	memberDeal := NewMemberDeal(message.GroupId, message.SelfId, false)
	list, err := memberDeal.GetMemberInfoList()
	if err != nil {
		logger.Log.Errorf("[GetMemberInfoList failed][err:%v]", err)
		return
	}

	MemberList[message.GroupId] = &list.Data
	//todo: @全体成员

	/*
		先在MemberList中查询信息，如果没有 则访问api请求查询
	*/

	var qqList []int64
	for _, v := range list.Data {
		if v.UserId != message.SelfId {
			qqList = append(qqList, v.UserId)
		}
	}

	atAllstring := tools.GetAtAllMemberString(qqList)

	sendMsg := tools.GetGroupMsgStruct(atAllstring, message.GroupId)
	tools.MsgChan <- sendMsg
	logger.Log.Infof("[user %s use AtAllMember success]", message.Sender.NickName)
}
