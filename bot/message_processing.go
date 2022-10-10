package bot

import (
	"fmt"
	"log"
	"qbot/bot/at_member"
	"qbot/bot/common/tools"
	"qbot/bot/member_change_notice"
	"qbot/bot/weather"
	"qbot/pkg/utils"
)

const (
	cqMe  = "[CQ:at,qq=%s] "
	AtAll = "/艾特全体"

	PostMessage   = "message"
	PostRequest   = "request"
	PostNotice    = "notice"
	PostMetaEvent = "meta_event"
)

var (
	AtMe = fmt.Sprintf(cqMe, utils.GlobalConf.QqBot.QqId)
)

func MessageDistribution(messageBytes []byte) {

	message, err := tools.ParsingMessage(messageBytes)
	if err != nil {
		log.Println("parse message failed,message = ", string(messageBytes))
		return
	}

	switch message.PostType {
	//处理消息事件
	case PostMessage:
		//at机器人自己
		if message.Message == AtMe {
			go at_member.AtMeFunc(message)
		}

		//at全体成员
		if message.Message == AtAll {
			go at_member.AtAllMember(message)
		}

		//天气查询
		go weather.WeatherQueryFunc(message)

	//处理请求事件
	case PostRequest:
		fmt.Println("request")

	//处理通知事件
	case PostNotice:
		//群成员增加欢迎通知
		go member_change_notice.MemberJoinNotice(message)

	//处理元事件
	case PostMetaEvent:
		fmt.Println("meta_event")
	}

}
