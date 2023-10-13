package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
	"github.com/jzyong/ugk/agent-manager/manager"
	"github.com/jzyong/ugk/agent-manager/rpc"
	"runtime"
)

// agent manager 入口
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfigAndLog()
	log.Info("启动 agent manager ......")

	var err error
	err = m.Init()
	if err != nil {
		log.Error("agent manager 启动错误: %s", err.Error())
		return
	}

	m.Run()
	util.WaitForTerminate()
	m.Stop()

}

type ModuleManager struct {
	*util.DefaultModuleManager
	AgentManager  *manager.AgentManager
	DockerManager *manager.DockerManager
	GrpcManager   *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.AgentManager = m.AppendModule(manager.GetAgentManager()).(*manager.AgentManager)
	m.DockerManager = m.AppendModule(manager.GetDockerManager()).(*manager.DockerManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
