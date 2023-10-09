package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
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
}

func (m *DockerManager) Stop() {
}
