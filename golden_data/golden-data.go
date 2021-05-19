package main

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RepoInfo struct {
	reponame string
	url      string
	branch   string
	database string
}

//非阻塞式执行命令
func ExecCmdNoWait(command string) string {
	cmd := exec.Command("/bin/bash", "-c", command)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Error(stderr.String())
	} else {
		log.Info(out.String())
	}
	resp := out.String()
	return resp
}

//阻塞式执行命令
func ExecCmdWait(command string) bool {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
		return false
	}
	cmd.Start()

	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
	return true
}

//获取时间
func GetTime() (string, string) {
	curretDate := time.Now()
	d, _ := time.ParseDuration("-24h")
	beforeData := curretDate.Add(d).Format("2006-01-02")
	return curretDate.Format("2006-01-02"), beforeData
}

func DowloadBeforeOSS() {
	curretDate, beforeData := GetTime()
	//判断目录是否存在
	if _, err := os.Stat(curretDate); os.IsNotExist(err) {
		command := "ossutil cp -r oss://ep-gold-data/" + beforeData + " ." + " && " + "mv " + beforeData + " " + curretDate
		log.Info(command)
		_ = ExecCmdWait(command)
	} else {
		command := "rm -rf " + curretDate + " && " + "ossutil cp -r oss://ep-gold-data/" + beforeData + " ." + " && " + "mv " + beforeData + " " + curretDate
		log.Info(command)
		_ = ExecCmdWait(command)
	}
}

func UploadCurretOSS() {
	curretDate, _ := GetTime()
	//判断目录是否存在
	command := "echo 'y'|ossutil rm -r oss://ep-gold-data/" + curretDate + " && ossutil cp -r " + curretDate + "/" + " oss://ep-gold-data/" + curretDate
	log.Info(command)
	_ = ExecCmdWait(command)
}

func GetCommitSqlFilePath(repoinfo RepoInfo, curretDate string) []string {
	//拉取代码
	if _, err := os.Stat(repoinfo.reponame); os.IsNotExist(err) {
		command := "git clone -b " + repoinfo.branch + " " + repoinfo.url
		_ = ExecCmdNoWait(command)
	} else {
		command := "cd " + repoinfo.reponame + " && git pull"
		_ = ExecCmdNoWait(command)
	}
	//获取相对时间的commit sql path
	command := "cd " + repoinfo.reponame + " && git log --after={" + curretDate + "} -p|grep '^+++'|grep 'sql'|awk '{print $NF}'|sed 's#b/##'|sort|uniq"
	//fmt.Println(command)
	resp := ExecCmdNoWait(command)
	paths := strings.Split(resp, "\n")
	return paths
}

func MergeCommitSqlFile(repoinfo RepoInfo, paths []string) {
	curretDate, _ := GetTime()
	baseDir, _ := os.Getwd()
	commitFile := baseDir + "/" + curretDate + "/" + repoinfo.database + "/" + "update-" + curretDate + ".sql"
	//因为是追加文件，防止多次跑
	//if _, err := os.Stat(commitFile); os.IsNotExist(err) {

	//} else {
	//	command := "rm -rf " + commitFile
	//	_ = ExecCmdNoWait(command)
	//}
	//却掉因为split切割导致的最后路径为空，报错
	for i := 0; i < len(paths)-1; i++ {
		command := "cat " + baseDir + "/" + repoinfo.reponame + "/" + paths[i] + " >> " + commitFile
		_ = ExecCmdNoWait(command)
	}
}

//数据集导入
type InstanceInfo struct {
	env      string
	label    string
	version  string
	host     string
	port     string
	username string
	password string
	dbs      []string
}

func GetDirFile(path string) []string {
	var temp []string
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error(err)
	}
	for i := range fileInfoList {
		//fmt.Println(fileInfoList[i].Name())
		temp = append(temp, fileInfoList[i].Name())
	}
	return temp
}

func InputInstanceFromFile(instances []InstanceInfo) {
	curretDate, beforeData := GetTime()
	//判断目录是否存在
	if _, err := os.Stat(curretDate); os.IsNotExist(err) {
		command := "ossutil cp -r oss://ep-gold-data/" + beforeData + " ." + " && " + "mv " + beforeData + " " + curretDate
		log.Info(command)
		_ = ExecCmdWait(command)
	}
	ossdbs := GetDirFile(curretDate)
	for _, ossdb := range ossdbs {
		for _, instance := range instances {
			for _, instancedb := range instance.dbs {
				if ossdb == instancedb {
					sqlfiles := GetDirFile(curretDate + "/" + ossdb)
					for _, sqlfile := range sqlfiles {
						filepath := curretDate + "/" + ossdb + "/" + sqlfile
						if instance.version == "mysql56" || instance.version == "mysql57" {
							command := "mysql -h" + instance.host + " -P" + instance.port + " -u" + instance.username + " -p'" + instance.password + "' " + instancedb + "_" + instance.env + " < " + filepath
							log.Info(command)
							_ = ExecCmdWait(command)
						} else {
							log.Info(ossdb, "pgsql10")
						}
					}
				}
			}
		}
	}

}

func init() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.SetOutput(os.Stdout)
}

func main() {
	//获取数据集
	repoinfos := []RepoInfo{
		{
			reponame: "achievement-service",
			url:      "git@github.com:WeWork-China/achievement-service.git",
			branch:   "master",
			database: "achievement",
		},
		{
			reponame: "china-building-info-service",
			url:      "git@github.com:WeWork-China/china-building-info-service.git",
			branch:   "master",
			database: "mulan_bis",
		},
		{
			reponame: "chinaos-faceservice",
			url:      "git@github.com:WeWork-China/chinaos-faceservice.git",
			branch:   "master",
			database: "wwc_face",
		},
	}

	curretDate, _ := GetTime()
	DowloadBeforeOSS()
	for _, repoinfo := range repoinfos {
		paths := GetCommitSqlFilePath(repoinfo, curretDate)
		MergeCommitSqlFile(repoinfo, paths)
		log.Info(repoinfo)
	}
	UploadCurretOSS()

	//数据集导入
	instanceinfos := []InstanceInfo{
		{
			env:      "test",
			label:    "mysql-java-test-mulan-db-v56",
			version:  "mysql56",
			host:     "rm-uf6g6uhcktm12y1ik5o.mysql.rds.aliyuncs.com",
			port:     "3306",
			username: "dms",
			password: "q0OzU4B^bpnEqkS4",
			dbs:      []string{"wwc_face"},
		},
		{
			env:      "test",
			label:    "mysql-java-test-mulan-db-v57",
			version:  "mysql57",
			host:     "rm-uf65m74z317qid2i8oo.mysql.rds.aliyuncs.com",
			port:     "3306",
			username: "dms",
			password: "q0OzU4B^bpnEqkS4",
			dbs:      []string{"mulan_bis"},
		},
	}
	InputInstanceFromFile(instanceinfos)
}
