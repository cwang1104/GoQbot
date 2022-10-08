package ws

import (
	"fmt"
	"qbot/bot"
	"qbot/bot/at_member"
	"qbot/bot/common/tools"
	"qbot/bot/weather"
	"qbot/pkg/utils"
)

const (
	cqMe  = "[CQ:at,qq=%s] "
	AtAll = "/艾特全体"
)

var (
	MsgChan    = make(chan *tools.SendWsMsgModel, 100)
	MemberList = map[int64]*[]bot.MemberInfo{} //群成员信息，用于at全体成员，加载完成后只读，不涉及写，不用锁

	AtMe = fmt.Sprintf(cqMe, utils.GlobalConf.QqBot.QqId)
)

func MessageDistribution(message *MessageType) {

	//@me 功能列表
	if message.Message == AtMe {
		go at_member.AtMeFunc(message)
	}

	if message.Message == AtAll {
		go at_member.AtAllMember(message)
	}

	//天气查询
	go weather.WeatherQueryFunc(message)
}
