package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"log"
	"os"
	"strconv"
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

func GetUserList() {
	request := ram.CreateListUsersRequest()
	request.Scheme = "https"

	response, err := client.ListUsers(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	policies := make(map[string]string)
	//fmt.Printf("response is %#v\n", response.Users.User)
	for _, v := range response.Users.User {
		//fmt.Println(GetUserPolicy(v.UserName))
		policies[v.UserName] = GetUserPolicy(v.UserName)

	}
	//for k,v := range policies{
	//	fmt.Println(k,v)
	//}
	WriteExcel(policies)
}

func GetGroupList() {
	request := ram.CreateListGroupsRequest()
	request.Scheme = "https"

	response, err := client.ListGroups(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	policies := make(map[string]string)
	//fmt.Printf("response is %#v\n", response.Users.User)
	for _, v := range response.Groups.Group {
		//fmt.Println(GetUserPolicy(v.UserName))
		policies[v.GroupName] = GetGroupPolicy(v.GroupName)

	}
	for k, v := range policies {
		fmt.Println(k, v)
	}
	WriteExcel(policies)
}

func GetUserPolicy(username string) string {
	request := ram.CreateListPoliciesForUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, err := client.ListPoliciesForUser(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response.Policies)
	policies := ""
	for _, v := range response.Policies.Policy {
		//policies = append(policies, v.PolicyName)
		//policies = append(policies, v.PolicyName)
		policies += v.PolicyName + ","
	}
	return policies
}

func GetGroupPolicy(groupname string) string {
	request := ram.CreateListPoliciesForGroupRequest()
	request.Scheme = "https"
	request.GroupName = groupname
	response, err := client.ListPoliciesForGroup(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response.Policies)
	policies := ""
	for _, v := range response.Policies.Policy {
		//policies = append(policies, v.PolicyName)
		//policies = append(policies, v.PolicyName)
		policies += v.PolicyName + ","
	}
	return policies
}

func WriteExcel(policies map[string]string) {
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("userpolicy")
	i := 1
	for k, v := range policies {
		//设置单元格的值
		fmt.Println("A"+strconv.Itoa(i), k)
		fmt.Println("B"+strconv.Itoa(i), v)
		xlsx.SetCellValue("userpolicy", "A"+strconv.Itoa(i), k)
		xlsx.SetCellValue("userpolicy", "B"+strconv.Itoa(i), v)
		i++
	}
	xlsx.SetActiveSheet(index)
	err := xlsx.SaveAs("./policy.xlsx")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Init()
	GetUserList()
	//GetGroupList()
}
