package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"os"
	//"strconv"
)

var client *rds.Client
var err error

var (
	region string
	key    string
	secret string
	action string
	tag    string
)

func Init() {
	flag.StringVar(&region, "r", "cn-shanghai", "region")
	flag.StringVar(&key, "k", "", "accessKeyId")
	flag.StringVar(&secret, "s", "", "accessSecret")
	flag.StringVar(&action, "a", "", "action：createbackup、")
	flag.StringVar(&tag, "t", "", "tag")
	flag.Parse()
	client, err = rds.NewClientWithAccessKey(region, key, secret)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
func CreateBackupJob(instanceid string) string {
	request := rds.CreateCreateBackupRequest()
	request.Scheme = "https"
	request.DBInstanceId = instanceid
	request.BackupStrategy = "instance"
	//根据磁盘类型配置备份方式
	disktype := rds.CreateDescribeDBInstanceAttributeRequest()
	disktype.Scheme = "https"
	disktype.DBInstanceId = instanceid
	disktyperesponse, err := client.DescribeDBInstanceAttribute(disktype)
	if err != nil {
		fmt.Print(err.Error())
	}
	if disktyperesponse.Items.DBInstanceAttribute[0].DBInstanceStorageType == "local_ssd" {
		request.BackupMethod = "Physical"
	} else {
		request.BackupMethod = "Snapshot"
	}

	request.BackupType = "FullBackup"
	response, err := client.CreateBackup(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
	return response.BackupJobId
}

func DescribeBackupTasks(instanceid string, backupjobid int) string {
	request := rds.CreateDescribeBackupTasksRequest()
	request.Scheme = "https"

	request.DBInstanceId = instanceid
	request.BackupJobId = requests.NewInteger(backupjobid)

	response, err := client.DescribeBackupTasks(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)
	return response.Items.BackupJob[0].BackupId
}

func DescribeBackups(instanceid string, backupid string) {
	request := rds.CreateDescribeBackupsRequest()
	request.Scheme = "https"

	request.DBInstanceId = instanceid
	request.BackupId = backupid

	response, err := client.DescribeBackups(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}

func DescribeDBInstanceAttribute(instanceid string) {
	request := rds.CreateDescribeDBInstanceAttributeRequest()
	request.Scheme = "https"

	request.DBInstanceId = "rm-uf65piy3508urmlj7"

	response, err := client.DescribeDBInstanceAttribute(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}

func CloneDBInstance(instanceid, backupid string) {
	request := rds.CreateCloneDBInstanceRequest()
	request.Scheme = "https"

	request.DBInstanceId = instanceid
	request.DBInstanceStorageType = "local_ssd"
	request.PayType = "Postpaid"
	request.BackupId = backupid

	response, err := client.CloneDBInstance(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}

func main() {
	Init()
	if action == "createbackup" {
		instanceid := "pgm-uf6sjq9qdalx0wi7"
		//backupjobid,_:= strconv.Atoi(CreateBackupJob(instanceid))
		backupfileid := DescribeBackupTasks(instanceid, 12618910)
		fmt.Println(instanceid, tag, backupfileid)
	}
	//DescribeBackupTasks("rm-uf65piy3508urmlj7",12544048)
	//DescribeBackups("rm-uf65piy3508urmlj7","832201883")
	//DescribeDBInstanceAttribute("rm-uf65piy3508urmlj7")
	//CloneDBInstance("rm-uf65piy3508urmlj7","832201883")
}
