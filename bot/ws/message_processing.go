package ws

import (
	"fmt"
	"log"
	"qbot/bot"
	"qbot/bot/http"
	"qbot/utils"
	"sync"
)

const (
	AtMe         = "[CQ:at,qq=903086461] "
	AtAll        = "/艾特全体"
	weatherQuery = "/天气"
)

var (
	weatherUserList = map[int64]map[int64]int64{} //群号对应下的QQ号列表
	mux             sync.Mutex
	MsgChan         = make(chan *SendWsMsgModel, 100)
	MemberList      = map[int64]*[]bot.MemberInfo{} //加载完成后只读，不涉及写，不用锁
	supportCity     string
)

func init() {
	supportCity = "天气查询支持以下城市:\n"
	for k, _ := range utils.GlobalConf.QqBot.WeatherLocation {
		supportCity = supportCity + k + "\n"
	}
	supportCity = supportCity + "\n请直接输入要查询的城市名称\n 输入 退出 退出天气查询"
}

func MessageDistribution(message *MessageType) {

	//@me 功能列表
	if message.Message == AtMe {
		go AtMeFunc(message)
	}

	if message.Message == AtAll {
		go AtAllMember(message)
	}

	//天气查询
	go WeatherQueryFunc(message)

}

func AtMeFunc(message *MessageType) {
	msg := "功能菜单：\n 1-查询天气 (输入 /天气)\n 2-艾特全体 (输入 /艾特全体)\n"
	sendInfo := NewSendGroupMsg(msg, message.GroupId)
	MsgChan <- sendInfo
}

//WeatherQueryFunc 天气查询功能实现函数
//传入参数为 MessageType指针
func WeatherQueryFunc(message *MessageType) {

	if message.Message == weatherQuery {
		//是否当前已经在天气查询中
		_, exist := GetWeatherUser(message.GroupId, message.UserId)
		if !exist {
			go func() {
				AddWeatherUser(message.GroupId, message.UserId)
				sendMsg := NewSendGroupMsg(supportCity, message.GroupId)
				MsgChan <- sendMsg
			}()
			return
		}
	}

	//获取当前已经在天气查询中的qq号
	userList := GetWeatherUserList(message.GroupId)
	fmt.Println("-----------", userList)
	for _, v := range userList {
		if message.UserId == v {
			if citySupport(message.Message) {
				cityCode := utils.GlobalConf.QqBot.WeatherLocation[message.Message]
				weatherString, err := utils.NewWeatherProvider(cityCode).GetWeatherString()
				if err != nil {
					log.Println("get weather info failed" + err.Error())
					sendMsg := NewSendGroupMsg("获取天气失败", message.GroupId)
					MsgChan <- sendMsg
					return
				}
				sendMsg := NewSendGroupMsg(weatherString, message.GroupId)
				MsgChan <- sendMsg
			} else if message.Message == "退出" {
				fmt.Println("当前输入：", message.Message)
				DelWeatherUser(message.GroupId, message.UserId)
				sendMsg := NewSendGroupMsg("好的，退出", message.GroupId)
				MsgChan <- sendMsg
			} else {
				sendMsg := NewSendGroupMsg(supportCity, message.GroupId)
				MsgChan <- sendMsg
			}
		}
	}
}

// AtAllMember @全体成员
func AtAllMember(message *MessageType) {

	if message.Sender.Role == "member" {
		sendMsg := NewSendGroupMsg("此功能仅群主及管理员可用", message.GroupId)
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

	sendMsg := NewSendGroupMsg(atAllstring, message.GroupId)
	MsgChan <- sendMsg

}

//NewSendGroupMsg 获取发送qq群消息内容
func NewSendGroupMsg(message string, groupId int64) *SendWsMsgModel {
	msg := SendGroupMsg{
		GroupID:    groupId,
		Message:    message,
		AutoEscape: false,
	}

	wsMsg := SendWsMsgModel{
		Action: "send_group_msg",
		Params: msg,
	}
	return &wsMsg
}

func AddWeatherUser(groupId, userId int64) {
	mux.Lock()
	fmt.Println("--1", weatherUserList)
	defer mux.Unlock()
	_, exist := weatherUserList[groupId]
	fmt.Println("--2", weatherUserList)
	if !exist {
		fmt.Println("--3", weatherUserList)
		userIdMap := make(map[int64]int64)
		userIdMap[userId] = userId
		weatherUserList[groupId] = userIdMap
		fmt.Println("--4", weatherUserList)
		return
	}
	fmt.Println("--5")
	weatherUserList[groupId][userId] = userId
	fmt.Println("--6")
}

func GetWeatherUserList(groupId int64) []int64 {
	mux.Lock()
	var list []int64
	for _, v := range weatherUserList[groupId] {
		list = append(list, v)
	}
	mux.Unlock()
	return list
}

func GetWeatherUser(groupId, userId int64) (id int64, exist bool) {
	mux.Lock()
	defer mux.Unlock()
	gid, exist := weatherUserList[groupId]
	if exist {
		id, exist = gid[userId]
	}
	return
}

func DelWeatherUser(groupId, userId int64) {
	mux.Lock()
	delete(weatherUserList[groupId], userId)
	mux.Unlock()
}

//locationSupport 判断输入的城市是否支持
func citySupport(city string) (exist bool) {
	_, exist = utils.GlobalConf.QqBot.WeatherLocation[city]
	return
}
