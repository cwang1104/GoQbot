package bothttp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"net/http"
)

var (
	jobList = map[string]*Job{}
)

type Job struct {
	cro     *cron.Cron
	status  int
	JobName string
}

func JobList(c *gin.Context) {
	var NameList []string
	for _, v := range jobList {
		fmt.Println("-------", v)
		if v.status == 1 {
			NameList = append(NameList, v.JobName)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"list": NameList,
	})
}

func DelJob(c *gin.Context) {
	name := "test"
	cro := jobList[name]
	cro.cro.Stop()
	cro.status = 2
	fmt.Println("ok")
}

func AddDsq(c *gin.Context) {
	name := "test"
	Addcro(name)
}

func Addcro(name string) {
	c := cron.New()

	spec := "*/5 * * * * ?"
	i := 1
	c.AddFunc(spec, func() {
		i++
		fmt.Println("list", i)
	})
	c.Start()
	testjob := &Job{
		cro: c,
	}
	testjob.JobName = name
	jobList[name] = testjob
}
