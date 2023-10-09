package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"os/exec"
	"runtime"
	"testing"
)

// 执行系统命令
func TestCmd(t *testing.T) {
	if runtime.GOOS == "windows" {
		windowCmd()
	} else if runtime.GOOS == "linux" {
		linuxShell()
	}
}

// window命令
func windowCmd() {
	cmd := exec.Command("docker", "ps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
		return
	}
	log.Info("执行结果：%v", string(output))
}

// linux命令
func linuxShell() {
	cmd := exec.Command("docker", "ps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
}

// 启动容器
func TestStartDocker(t *testing.T) {
	//1.关闭
	//docker -H tcp://127.0.0.1:2375 stop slots-api-jzy1
	cmd := exec.Command("docker", "-H", "tcp://127.0.0.1:2375", "stop", "slots-api-jzy1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
	}
	log.Info("关闭结果：%v", string(output))

	//2.删除
	//docker -H tcp://127.0.0.1:2375 rm  slots-api-jzy1
	cmd = exec.Command("docker", "-H", "tcp://127.0.0.1:2375", "rm", "slots-api-jzy1")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
	}
	log.Info("删除结果：%v", string(output))

	//3.运行
	//set javaParam="-server -Xmx256M -Xms256M"
	//docker -H tcp://127.0.0.1:2375 run -dit -p 7060:7060 --name slots-api-jzy1 --restart=always -e JAVA_OPTS=%javaParam% 127.0.0.1:5000/jzy1/slots-api:releases
	cmd = exec.Command("docker", "-H", "tcp://127.0.0.1:2375", "run", "-dit", "-p", "7060:7060", "--name", "slots-api-jzy1", "--restart=always", "-e", "JAVA_OPTS=-server -Xmx256M -Xms256M", "slots-api:releases")
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
	}
	log.Info("运行结果：%v", string(output))
}
