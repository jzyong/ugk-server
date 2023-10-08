package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
	"github.com/jzyong/ugk/agent-manager/manager"
	"github.com/jzyong/ugk/agent-manager/rpc"
	"runtime"
)

// gate 入口 TODO grpc服务添加
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
	AgentManager *manager.AgentManager
	GrpcManager  *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.AgentManager = m.AppendModule(manager.GetAgentManager()).(*manager.AgentManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
