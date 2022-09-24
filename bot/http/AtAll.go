package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"qbot/bot"
)

const (
	baseUrl = "http://127.0.0.1:5700"
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

func (g *GroupMemberDeal) AtAllMember() error {

	//获取群成员列表
	memberList, err := g.GetMemberInfoList()
	if err != nil {
		log.Println("GetMemberInfoList failed", err)
		return err
	}

	//获取群成员qq号列表
	var memberIdList []int64
	for _, v := range memberList.Data {
		if v.Role != "owner" {
			memberIdList = append(memberIdList, v.UserId)
		}
		//memberIdList = append(memberIdList, v.UserId)
	}

	//type Msg struct {
	//	Type string `json:"type"`
	//	Data struct {
	//		QQ int64 `json:"qq"`
	//	} `json:"data"`
	//}
	//var msg Msg
	//msg.Type = "CQ"
	//msg.Data.QQ = 45241397
	//
	////data, _ := json.Marshal(msg)
	////message := []string{string(data)}

	//message := "[cq:at"
	//for _, v := range memberIdList {
	//	message += fmt.Sprintf(",qq=%d", v)
	//}
	//message = message + "]"
	//message := fmt.Sprintf("[cq:at,qq=%d]", memberIdList[0])
	//message := "&#91CQ:face&#44id=1&#93"
	//log.Println("message = ", message)

	message2 := fmt.Sprintf("[CQ:at,qq=45241397]测试")
	err = g.SendGroupMsg(message2, false)
	if err != nil {
		log.Println("SendGroupMsg failed", err)
		return err
	}
	return nil
}

func (g *GroupMemberDeal) SendGroupMsg(msg string, autoEscape bool) error {

	sendMsg := bot.SendGroupMsg{
		GroupID:    g.GroupId,
		Message:    msg,
		AutoEscape: autoEscape,
	}
	b, err := json.Marshal(sendMsg)
	fmt.Println(string(b))
	if err != nil {
		log.Println("sendMsg json Marshal failed", err)
		return err
	}
	url := fmt.Sprintf("%s/send_group_msg", baseUrl)
	respR, err := httpPost(url, b)
	if err != nil {
		log.Println("httpPost failed", err)
		return err
	}

	log.Println("string-------", string(respR))

	//var memberResp MemberListResp
	//err = json.Unmarshal(respR, &memberResp)
	//if err != nil {
	//	log.Println("Unmarshal,response err", err)
	//	return  err
	//}
	return nil
}

func httpPost(url string, b []byte) ([]byte, error) {
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		log.Println("Post failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", err)
		return nil, err
	}
	return content, nil
}
