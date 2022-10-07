package ws

import (
	"fmt"
	"log"
	"qbot/bot"
	"qbot/bot/http"
	"qbot/bot/weather"
	"qbot/utils"
)

const (
	cqMe  = "[CQ:at,qq=%s] "
	AtAll = "/艾特全体"
)

var (
	MsgChan    = make(chan *SendWsMsgModel, 100)
	MemberList = map[int64]*[]bot.MemberInfo{} //群成员信息，用于at全体成员，加载完成后只读，不涉及写，不用锁

	AtMe = fmt.Sprintf(cqMe, utils.GlobalConf.QqBot.QqId)
)

func MessageDistribution(message *MessageType) {

	//@me 功能列表
	if message.Message == AtMe {
		go AtMeFunc(message)
	}

	if message.Message == AtAll {
		go AtAllMember(message)
	}

	//天气查询
	go weather.WeatherQueryFunc(message)

}

func AtMeFunc(message *MessageType) {
	msg := "功能菜单：\n 1-查询天气 (输入 /天气)\n 2-艾特全体 (输入 /艾特全体)\n"
	sendInfo := GetGroupMsgStruct(msg, message.GroupId)
	MsgChan <- sendInfo
}

// AtAllMember @全体成员
func AtAllMember(message *MessageType) {

	if message.Sender.Role == "member" {
		sendMsg := GetGroupMsgStruct("此功能仅群主及管理员可用", message.GroupId)
		MsgChan <- sendMsg
		return
	}

	//获取群成员的信息，构建qq号切片
	memberDeal := http.NewMemberDeal(message.GroupId, message.SelfId, false)
	list, err := memberDeal.GetMemberInfoList()
	if err != nil {
		log.Println("GetMemberInfoList failed", err)
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

	atAllstring := utils.GetAtAllMemberString(qqList)

	sendMsg := GetGroupMsgStruct(atAllstring, message.GroupId)
	MsgChan <- sendMsg

}

