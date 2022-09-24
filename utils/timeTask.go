package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GetTimeTaskName(taskName string, taskId int) string {
	return fmt.Sprintf("%s-%d", taskName, taskId)
}

//GetInternalSpec 根据internal 获取cron表达式
func GetInternalSpec(interval, timeLimitStart, timeLimitEnd int) string {
	//拼接cron表达式
	var spec string

	//时间间隔为15分钟或30分钟
	if interval == 15 || interval == 30 {
		spec = fmt.Sprintf("* */%d %d-%d * * ?", interval, timeLimitStart, timeLimitEnd)

		//时间间隔为1小时或者2小时
	} else if interval == 60 || interval == 120 {
		var intervalList []string
		var intervalString string
		if interval == 60 {
			for i := timeLimitStart; i <= timeLimitEnd; i++ {
				intervalString = strconv.Itoa(i)
				intervalList = append(intervalList, intervalString)
			}
		} else if interval == 120 {
			for i := timeLimitStart; i <= timeLimitEnd; i = i + 2 {
				intervalString = strconv.Itoa(i)
				intervalList = append(intervalList, intervalString)
			}
		}
		intervalListString := strings.Join(intervalList, ",")
		spec = fmt.Sprintf("* * %s * * ?", intervalListString)
	}
	return spec
}
