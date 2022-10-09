package at_member

import (
	"encoding/json"
	"fmt"
	"log"
	"qbot/bot/common/tools"
)

var MemberList = map[int64]*[]MemberInfo{} //群成员信息，用于at全体成员，加载完成后只读，不涉及写，不用锁
//MemberInfo 成员信息
type MemberInfo struct {
	GroupId         int64  `json:"group_id"`
	UserId          int64  `json:"user_id"`
	NickName        string `json:"nick_name"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int32  `json:"age"`
	Area            string `json:"area"`
	JoinTime        int32  `json:"join_time"`
	LastSentTime    int32  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int64  `json:"title_expire_time"`
	CardChangeAble  bool   `json:"card_change_able"`
	ShutUpTimestamp int64  `json:"shut_up_timestamp"`
}

//GetMemberList 获取群成员api请求参数
type GetMemberList struct {
	GroupId int64 `json:"group_id"`
	NoCache bool  `json:"no_cache"`
}

type GroupMemberDeal struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
	NoCache bool  `json:"no_cache"`
}

type MemberListResp struct {
	Data    []MemberInfo `json:"data"`
	RetCode int          `json:"retcode"`
	Status  string       `json:"status"`
}

func NewMemberDeal(groupId, userId int64, cache bool) *GroupMemberDeal {
	return &GroupMemberDeal{
		GroupId: groupId,
		UserId:  userId,
		NoCache: cache,
	}
}

func (g *GroupMemberDeal) GetMemberInfoList() (*MemberListResp, error) {

	sendMsg := GetMemberList{
		GroupId: g.GroupId,
		NoCache: g.NoCache,
	}
	b, err := json.Marshal(sendMsg)
	if err != nil {
		log.Println("sendMsg json Marshal failed", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/get_group_member_list", tools.CqHttpBaseUrl)
	respR, err := tools.HttpPost(url, b)
	if err != nil {
		log.Println("httpPost failed", err)
		return nil, err
	}

	var memberResp MemberListResp
	err = json.Unmarshal(respR, &memberResp)
	if err != nil {
		log.Println("Unmarshal,response err", err)
		return nil, err
	}
	return &memberResp, nil
}
