package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
		fmt.Println(stderr.String())
	} else {
		fmt.Println(out.String())
	}
	resp := out.String()
	return resp
}

//阻塞式执行命令
func ExecCmdWait(command string) bool {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(command)
		_ = ExecCmdWait(command)
	} else {
		command := "rm -rf " + curretDate + " && " + "ossutil cp -r oss://ep-gold-data/" + beforeData + " ." + " && " + "mv " + beforeData + " " + curretDate
		_ = ExecCmdWait(command)
	}
}

func UploadCurretOSS() {
	curretDate, _ := GetTime()
	//判断目录是否存在
	command := "echo 'y'|ossutil rm -r oss://ep-gold-data/" + curretDate + " && ossutil cp -r " + curretDate + "/" + " oss://ep-gold-data/" + curretDate
	fmt.Println(command)
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
	commitFile := baseDir + "/" + curretDate + "/" + repoinfo.database + "/" + "add-" + curretDate + ".sql"
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
func main() {
	//job功能
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
		fmt.Println(repoinfo)
	}
	UploadCurretOSS()

	//
}
