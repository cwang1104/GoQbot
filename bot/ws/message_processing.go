package ws

import (
	"fmt"
	"log"
	"qbot/bot"
	"qbot/bot/http"
	"sync"
)

const (
	WeatherStart = "请输入城市：\n1-[成都 眉山]\n2-输入”退出“ 退出查询"
	CityError    = "输入的城市错误，只能是 成都或眉山"
	WeatherResp  = "天气很好"

	AtMe         = "[CQ:at,qq=903086461] "
	weatherQuery = "/天气"
)

var (
	weatherUserList = map[int64]map[int64]int64{} //群号对应下的QQ号列表
	mux             sync.Mutex
	MsgChan         = make(chan *SendWsMsgModel, 100)
	MemberList      = map[int64]*[]bot.MemberInfo{}
)

func MessageDistribution(message *MessageType) {

	//@me 功能列表
	if message.Message == AtMe {
		go AtMeFunc(message)
	}

	//天气查询
	go WeatherQueryFunc(message)

}

func AtMeFunc(message *MessageType) {
	msg := "功能菜单：\n 1-查询天气 (输入 /天气)\n"
	sendInfo := NewSendGroupMsg(msg, message.GroupId)
	MsgChan <- sendInfo
}

func WeatherQueryFunc(message *MessageType) {

	if message.Message == weatherQuery {
		//是否当前已经在天气查询中
		_, exist := GetWeatherUser(message.GroupId, message.UserId)
		if !exist {
			go func() {
				AddWeatherUser(message.GroupId, message.UserId)
				sendMsg := NewSendGroupMsg(WeatherStart, message.GroupId)
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
			if message.Message == "成都" || message.Message == "眉山" {
				fmt.Println("当前输入：", message.Message)
				weatherInfo := fmt.Sprintf("%s%s", message.Message, WeatherResp)
				sendMsg := NewSendGroupMsg(weatherInfo, message.GroupId)
				MsgChan <- sendMsg
			} else if message.Message == "退出" {
				fmt.Println("当前输入：", message.Message)
				DelWeatherUser(message.GroupId, message.UserId)
				sendMsg := NewSendGroupMsg("好的，退出", message.GroupId)
				MsgChan <- sendMsg
			} else {
				fmt.Println("当前输入：", message.Message)
				sendMsg := NewSendGroupMsg(CityError, message.GroupId)
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

}

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
