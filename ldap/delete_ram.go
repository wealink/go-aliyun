package main

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/go-ldap/ldap/v3"
	"os"
	"strings"
)

var (
	client        *ram.Client
	err           error
	region        string
	key           string
	secret        string
	ldap_addr     string
	ldap_username string
	ldap_password string
)

func Init() {

	region = os.Getenv("REGION")
	key = os.Getenv("ACCESS_KEY")
	secret = os.Getenv("ACCESS_SECRET")
	ldap_addr = os.Getenv("LDAP_ADDR")
	ldap_username = os.Getenv("LDAP_USERNAME")
	ldap_password = os.Getenv("LDAP_PASSWORD")
	client, err = ram.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}

func ExitRamUser(username string) (code int) {
	request := ram.CreateGetUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.GetUser(request)
	return response.GetHttpStatus()
}

func DelteteMFA(username string) (code int) {
	request := ram.CreateUnbindMFADeviceRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.UnbindMFADevice(request)
	return response.GetHttpStatus()
}

func DelteteUserFromGroup(username string) {
	request := ram.CreateListGroupsForUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.ListGroupsForUser(request)
	for _, v := range response.Groups.Group {
		groupname := v.GroupName
		request1 := ram.CreateRemoveUserFromGroupRequest()
		request1.Scheme = "https"
		request1.UserName = username
		request1.GroupName = groupname
		client.RemoveUserFromGroup(request1)
	}
}

func DelteteUserPolicy(username string) {
	request := ram.CreateListPoliciesForUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.ListPoliciesForUser(request)
	for _, v := range response.Policies.Policy {
		//v.PolicyName
		request1 := ram.CreateDetachPolicyFromUserRequest()
		request1.Scheme = "https"
		request1.UserName = username
		request1.PolicyType = v.PolicyType
		request1.PolicyName = v.PolicyName
		client.DetachPolicyFromUser(request1)
	}
}

func GetUserKeyCount(username string) (count int) {
	request := ram.CreateListAccessKeysRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.ListAccessKeys(request)
	return len(response.AccessKeys.AccessKey)
}

func DelteteUserKey(username string) {
	request := ram.CreateListAccessKeysRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.ListAccessKeys(request)
	for _, v := range response.AccessKeys.AccessKey {
		request1 := ram.CreateDeleteAccessKeyRequest()
		request1.Scheme = "https"
		request1.UserName = username
		request1.UserAccessKeyId = v.AccessKeyId
		client.DeleteAccessKey(request1)
	}
}

func DeleteRamUser(username string) (code int) {
	request := ram.CreateDeleteUserRequest()
	request.Scheme = "https"
	request.UserName = username
	response, _ := client.DeleteUser(request)
	return response.GetHttpStatus()

}

func GetOffUserList() {
	conn, err := ldap.DialURL(ldap_addr)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	err = conn.Bind(ldap_username, ldap_password)
	if err != nil {
		fmt.Println(err)
	}

	searchRequest := ldap.NewSearchRequest("OU=Offboarding,DC=wework,DC=cn",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=*)",
		[]string{"mail"},
		nil)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(sr.Entries)
	if len(sr.Entries) > 0 {
		for _, item := range sr.Entries {
			mail := item.GetAttributeValue("mail")
			if mail != "" {
				username := strings.Split(mail, "@")[0]
				//fmt.Println(username)
				exitcode := ExitRamUser(username)
				if exitcode == 200 {
					fmt.Println("离职用户：" + username + "存在阿里云RAM")
					//获取access key个数
					keycount := GetUserKeyCount(username)
					if keycount == 0 {
						//解绑定MFA
						DelteteMFA(username)
						//将用户从组中删除
						DelteteUserFromGroup(username)
						//移除用户策略
						DelteteUserPolicy(username)
						//删除用户access key功能暂时不开启，防止业务代码中残留

						delcode := DeleteRamUser(username)
						if delcode == 200 {
							fmt.Println("删除离职用户：" + username + "阿里云RAM账号成功")
						} else {
							fmt.Println("删除离职用户："+username+"阿里云RAM账号失败", delcode)
						}
					} else {
						fmt.Println("离职用户：" + username + "存在Access Key，防止删除影响业务，跳过删除流程")
					}
				}
			}
		}
	}
}

func main() {
	Init()
	GetOffUserList()
}
