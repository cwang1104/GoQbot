package at_member

import (
	"encoding/json"
	"fmt"
	"log"
	"qbot/bot"
)

type GroupMemberDeal struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
	NoCache bool  `json:"no_cache"`
}

type MemberListResp struct {
	Data    []bot.MemberInfo `json:"data"`
	RetCode int              `json:"retcode"`
	Status  string           `json:"status"`
}

func NewMemberDeal(groupId, userId int64, cache bool) *GroupMemberDeal {
	return &GroupMemberDeal{
		GroupId: groupId,
		UserId:  userId,
		NoCache: cache,
	}
}

func (g *GroupMemberDeal) GetMemberInfoList() (*MemberListResp, error) {

	sendMsg := bot.GetMemberList{
		GroupId: g.GroupId,
		NoCache: g.NoCache,
	}
	b, err := json.Marshal(sendMsg)
	if err != nil {
		log.Println("sendMsg json Marshal failed", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/get_group_member_list", baseUrl)
	respR, err := httpPost(url, b)
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
