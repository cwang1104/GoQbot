package weather

import (
	"qbot/bot/common/tools"
	"qbot/pkg/logger"
	"qbot/pkg/utils"
)

const weatherQuery = "/天气"

//WeatherQueryFunc 天气查询功能实现函数
//传入参数为 MessageType指针
func WeatherQueryFunc(message *tools.MessageType) {
	if message.Message == weatherQuery {
		//是否当前已经在天气查询中
		_, exist := GetWeatherUser(message.GroupId, message.UserId)
		if !exist {
			go func() {
				AddWeatherUser(message.GroupId, message.UserId)
				sendMsg := tools.GetGroupMsgStruct(supportCity, message.GroupId)
				tools.MsgChan <- sendMsg
			}()
			return
		}
	}
	//获取当前已经在天气查询中的qq号
	userList := GetWeatherUserList(message.GroupId)

	for _, v := range userList {
		if message.UserId == v {
			logger.Log.Debugf("天气查询用户列表：%v", userList)
			if citySupport(message.Message) {
				cityCode := utils.GlobalConf.QqBot.WeatherLocation[message.Message]
				weatherString, err := tools.NewWeatherProvider(cityCode).GetWeatherString()
				if err != nil {
					logger.Log.Errorf("get weather failed;err:%v", err)
					sendMsg := tools.GetGroupMsgStruct("获取天气失败", message.GroupId)
					tools.MsgChan <- sendMsg
					return
				}
				sendMsg := tools.GetGroupMsgStruct(weatherString, message.GroupId)
				tools.MsgChan <- sendMsg
			} else if message.Message == "退出" {
				logger.Log.Infof("当前输入内容%v", message.Message)
				DelWeatherUser(message.GroupId, message.UserId)
				sendMsg := tools.GetGroupMsgStruct("好的，退出", message.GroupId)
				tools.MsgChan <- sendMsg
			} else {
				sendMsg := tools.GetGroupMsgStruct(supportCity, message.GroupId)
				tools.MsgChan <- sendMsg
			}
		}
	}
}
