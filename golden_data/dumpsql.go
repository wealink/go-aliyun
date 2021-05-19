package main

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type InstanceInfo struct {
	label    string
	version  string
	host     string
	port     string
	username string
	password string
	dbs      []string
}

//阻塞式执行命令
func ExecCmd(command string) bool {
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

func DumpDBFile(info InstanceInfo) {
	for _, db := range info.dbs {
		command := "mysqldump --column-statistics=0 -h" + info.host + " -p" + info.port + " -u" + info.username + " -p" + info.password + " " + db + " --extended-insert=true" + " > " + db + ".sql"
		fmt.Println(command)
		_ = ExecCmd(command)
	}
}

func main() {
	instances := []InstanceInfo{
		{
			label:    "mysql-java-prod-mulan-db-v57",
			version:  "mysql57",
			host:     "rm-uf65zox8h1mi44j67so.mysql.rds.aliyuncs.com",
			port:     "3306",
			username: "dms",
			password: "*v*#lgiqMuEKjheW",
			dbs:      []string{"mulan_bis_production"},
		},
		{
			label:    "mysql-java-prod-mulan-db-v56",
			version:  "mysql56",
			host:     "rm-uf6prcst021us4q3pfo.mysql.rds.aliyuncs.com",
			port:     "3306",
			username: "dms",
			password: "*v*#lgiqMuEKjheW",
			dbs:      []string{"wwc_face_production"},
		},
	}
	for _, v := range instances {
		DumpDBFile(v)
	}
}
