package at_member

import (
	"log"
	"qbot/bot/common/tools"
	"qbot/bot/http"
	"qbot/bot/ws"
)

func AtMeFunc(message *ws.MessageType) {
	msg := "功能菜单：\n 1-查询天气 (输入 /天气)\n 2-艾特全体 (输入 /艾特全体)\n"
	sendInfo := tools.GetGroupMsgStruct(msg, message.GroupId)
	ws.MsgChan <- sendInfo
}

// AtAllMember @全体成员
func AtAllMember(message *ws.MessageType) {

	if message.Sender.Role == "member" {
		sendMsg := tools.GetGroupMsgStruct("此功能仅群主及管理员可用", message.GroupId)
		ws.MsgChan <- sendMsg
		return
	}

	//获取群成员的信息，构建qq号切片
	memberDeal := http.NewMemberDeal(message.GroupId, message.SelfId, false)
	list, err := memberDeal.GetMemberInfoList()
	if err != nil {
		log.Println("GetMemberInfoList failed", err)
		return
	}

	ws.MemberList[message.GroupId] = &list.Data
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
	ws.MsgChan <- sendMsg
}
