package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RepoInfo struct {
	name   string
	url    string
	branch string
}

func ExecCmd(command string) string {
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

func GetCommitSqlFilePath(reponame, url, branch, afterdata string) []string {
	//拉取代码
	if _, err := os.Stat(reponame); os.IsNotExist(err) {
		command := "git clone -b " + branch + " " + url
		_ = ExecCmd(command)
	} else {
		command := "cd " + reponame + " && git pull"
		_ = ExecCmd(command)
	}
	//获取相对时间的commit sql path
	command := "cd " + reponame + " && git log --after={" + afterdata + "} -p|grep '^+++'|grep 'sql'|awk '{print $NF}'|sed 's#b/##'|sort|uniq"
	//fmt.Println(command)
	resp := ExecCmd(command)
	paths := strings.Split(resp, "\n")
	return paths
}

func MergeCommitSqlFile(reponame string, paths []string) {
	timeStr := time.Now().Format("2006-01-02")
	baseDir, _ := os.Getwd()
	commitFile := baseDir + "/" + reponame + "-" + timeStr + ".sql"
	//因为是追加文件，防止多次跑
	if _, err := os.Stat(commitFile); os.IsNotExist(err) {

	} else {
		command := "rm -rf " + commitFile
		_ = ExecCmd(command)
	}
	//却掉因为split切割导致的最后路径为空，报错
	for i := 0; i < len(paths)-1; i++ {
		command := "cat " + baseDir + "/" + reponame + "/" + paths[i] + " >> " + commitFile
		_ = ExecCmd(command)
	}
}
func main() {
	repoinfos := []RepoInfo{
		{
			name:   "china-member-service",
			url:    "git@github.com:WeWork-China/china-member-service.git",
			branch: "release",
		},
		{
			name:   "china-payment-service",
			url:    "git@github.com:WeWork-China/china-payment-service.git",
			branch: "master",
		},
		{
			name:   "mulan-credits-service",
			url:    "git@github.com:WeWork-China/mulan-credits-service.git",
			branch: "release",
		},
	}
	for _, repoinfo := range repoinfos {
		paths := GetCommitSqlFilePath(repoinfo.name, repoinfo.url, repoinfo.branch, "2021-03-29")
		MergeCommitSqlFile(repoinfo.name, paths)
	}
}
