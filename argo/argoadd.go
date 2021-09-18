package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Info struct {
	ClusterName string
	NameSpace   string
	Branch      string
	Env         string
}

//阻塞式执行命令
func ExecCmdWait(command string) bool {
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd.Start()
	BufferRead(stdout)
	BufferRead(stderr)

	//阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	cmd.Wait()
	return true
}

func BufferRead(message io.Reader) {
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(message)
	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		fmt.Println(line)
	}
}

func main() {
	serviceName := os.Args[1]
	infos := map[string]Info{
		"dev1": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev",
		},
		"dev2": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev",
		},
		"dev3": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev",
		},
		"dev4": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev",
		},
		"dev5": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev5",
			Branch:      "master",
			Env:         "dev5",
		},
		"dev6": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev6",
		},
		"dev7": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev7",
		},
		"dev8": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev8",
		},
		"dev9": Info{
			ClusterName: "enigma",
			NameSpace:   "china-tech-dev",
			Branch:      "master",
			Env:         "dev9",
		},
		"test": {
			ClusterName: "enigma",
			NameSpace:   "china-tech-test",
			Branch:      "master",
			Env:         "test",
		},

		"int": {
			ClusterName: "enigma",
			NameSpace:   "china-tech-int",
			Branch:      "auto-envs",
			Env:         "int",
		},
		"staging": {
			ClusterName: "oracle",
			NameSpace:   "china-tech-staging",
			Branch:      "auto-envs",
			Env:         "staging",
		},
		"production": {
			ClusterName: "oracle",
			NameSpace:   "china-tech-prod",
			Branch:      "auto-envs",
			Env:         "production",
		},
	}

	for env, info := range infos {
		cmd := fmt.Sprintf("argocd app create %s-%s --repo https://github.com/wework-china/china-self-service.git --path . --dest-name %s --dest-namespace %s --revision %s --project %s --config-management-plugin wework-cd --sync-policy auto ",
			serviceName, env, info.ClusterName, info.NameSpace, info.Branch, info.Env)
		fmt.Println(cmd)
		ExecCmdWait(cmd)
	}
}
