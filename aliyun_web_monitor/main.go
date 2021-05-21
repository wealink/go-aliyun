package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"os"
)

var client *cms.Client
var err error

func Init() {
	var (
		region string
		key    string
		secret string
	)
	region = os.Getenv("REGION")
	key = os.Getenv("ACCESS_KEY")
	secret = os.Getenv("ACCESS_SECRET")
	client, err = cms.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

type Item struct {
	TaskId   string
	TaskName string
}

func GetSiteMonitorList() []Item {
	request := cms.CreateDescribeSiteMonitorListRequest()
	request.Scheme = "https"

	response, err := client.DescribeSiteMonitorList(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	result := response.SiteMonitors.SiteMonitor
	items := make([]Item, len(result))
	for index, value := range result {
		items[index].TaskId = value.TaskId
		items[index].TaskName = value.TaskName
	}
	return items
}

func GetSiteMouthMetric(items []Item) {

}

func main() {
	Init()
	items := GetSiteMonitorList()
	fmt.Println(items)
}
