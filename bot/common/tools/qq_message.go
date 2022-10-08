package tools

import "fmt"

func GetAtAllMemberString(qqIdList []int64) string {
	var atAll string
	for _, v := range qqIdList {
		a := fmt.Sprintf("[CQ:at,qq=%d]", v)
		atAll = atAll + a
	}
	return atAll
}
