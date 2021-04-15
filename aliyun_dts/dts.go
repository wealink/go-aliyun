package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"os"
)

var client *ram.Client
var err error

func Init() {
	var (
		region string
		key    string
		secret string
	)
	flag.StringVar(&region, "r", "cn-shanghai", "region")
	flag.StringVar(&key, "k", "", "accessKeyId")
	flag.StringVar(&secret, "s", "", "accessSecret")
	flag.Parse()
	client, err = ram.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

func CreateDtsInstance() {
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dts.aliyuncs.com"
	request.Version = "2020-01-01"
	request.ApiName = "CreateDtsInstance"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["InstanceClass"] = "medium"
	request.QueryParams["PayType"] = "PostPaid"
	request.QueryParams["SyncArchitecture"] = "oneway"
	request.QueryParams["Quantity"] = "1"
	request.QueryParams["Type"] = "MIGRATION"
	request.QueryParams["SourceRegion"] = "cn-shanghai"
	request.QueryParams["DestinationRegion"] = "cn-shanghai"
	request.QueryParams["SourceEndpointEngineName"] = "MySQL"
	request.QueryParams["DestinationEndpointEngineName"] = "MySQL"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func main() {
	Init()
	CreateDtsInstance()
}
