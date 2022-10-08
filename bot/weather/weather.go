package weather

import (
	"qbot/pkg/utils"
	"sync"
)

var (
	//qq群中正在进行天气查询中的成员qq号列表
	weatherUserList = map[int64]map[int64]int64{} //群号对应下的QQ号列表
	mux             sync.Mutex

	//天气查询支持的城市
	supportCity string
)

func init() {
	supportCity = "天气查询支持以下城市:\n"
	for k, _ := range utils.GlobalConf.QqBot.WeatherLocation {
		supportCity = supportCity + k + "\n"
	}
	supportCity = supportCity + "\n请直接输入要查询的城市名称\n 输入 退出 退出天气查询"
}

func AddWeatherUser(groupId, userId int64) {
	mux.Lock()
	defer mux.Unlock()
	_, exist := weatherUserList[groupId]

	if !exist {
		userIdMap := make(map[int64]int64)
		userIdMap[userId] = userId
		weatherUserList[groupId] = userIdMap
		return
	}
	weatherUserList[groupId][userId] = userId
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
