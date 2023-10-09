package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent/config"
	"github.com/jzyong/ugk/agent/manager"
	"github.com/jzyong/ugk/agent/rpc"
	manager2 "github.com/jzyong/ugk/common/manager"
	"runtime"
)

// gate 入口 TODO grpc服务添加,连接agent-manager
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfigAndLog()
	log.Info("启动 gate ......")

	var err error
	err = m.Init()
	if err != nil {
		log.Error("gate 启动错误: %s", err.Error())
		return
	}

	m.Run()
	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()
}

type ModuleManager struct {
	*util.DefaultModuleManager
	AgentManager         *manager.AgentManager
	DockerManager        *manager.DockerManager
	MachineManager       *manager.MachineManager
	ServiceClientManager *manager2.ServiceClientManager
	GrpcManager          *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.AgentManager = m.AppendModule(manager.GetAgentManager()).(*manager.AgentManager)
	m.DockerManager = m.AppendModule(manager.GetDockerManager()).(*manager.DockerManager)
	m.MachineManager = m.AppendModule(manager.GetMachineManager()).(*manager.MachineManager)
	m.ServiceClientManager = m.AppendModule(manager2.GetServiceClientManager()).(*manager2.ServiceClientManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
