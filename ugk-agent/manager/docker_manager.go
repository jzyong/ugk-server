package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"os/exec"
	"sync"
)

// DockerManager 执行系统命令在docker中运行unity服务器
type DockerManager struct {
	util.DefaultModule
}

var dockerManager *DockerManager
var dockerSingletonOnce sync.Once

func GetDockerManager() *DockerManager {
	dockerSingletonOnce.Do(func() {
		dockerManager = &DockerManager{}
	})
	return dockerManager
}

func (m *DockerManager) Init() error {
	log.Info("DockerManager 初始化......")
	return nil
}
func (m *DockerManager) Run() {
	//TODO 临时测试
	go m.testRunCommand()
}

func (m *DockerManager) testRunCommand() {
	//执行docker ps
	cmd := exec.Command("docker", "ps")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("执行错误：%v", err)
		return
	}
	log.Info("执行结果：%v", string(output))
}

func (m *DockerManager) Stop() {
}
