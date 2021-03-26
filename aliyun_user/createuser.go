package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"os"
	"strings"
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

func CreateUser(username string) {
	request := ram.CreateCreateUserRequest()
	request.Scheme = "https"

	request.UserName = username
	request.DisplayName = username
	request.Email = username + "@wework.cn"

	response, err := client.CreateUser(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("创建用户：" + response.User.UserName + "成功")
}

func CreateLoginProfile(username, password string) {
	request := ram.CreateCreateLoginProfileRequest()
	request.Scheme = "https"

	request.PasswordResetRequired = requests.NewBoolean(true)
	request.Password = password
	request.MFABindRequired = requests.NewBoolean(true)
	request.UserName = username

	response, err := client.CreateLoginProfile(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println("创建用户登陆配置：" + response.LoginProfile.UserName + "成功")
}

func GetUserListFromExecl() []string {
	filename := "QuickBi用户分组"
	xlsx, err := excelize.OpenFile(filename + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
	var usernams = []string{}
	rows := xlsx.GetRows("Sheet1")
	for _, row := range rows {
		if row[2] != "账号" {
			usernams = append(usernams, row[2])
		}
	}
	return usernams
}

func main() {
	Init()
	for _, val := range GetUserListFromExecl() {
		username := strings.Split(val, "@")[0]
		CreateUser(username)
		CreateLoginProfile(username, "fT3O4WgxdjFI5yQ^")
	}
	//username := strings.Split("test@wework.cn","@")[0]
	//CreateUser(username)
	//CreateLoginProfile(username,"fT3O4WgxdjFI5yQ^")
}
