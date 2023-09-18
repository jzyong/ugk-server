package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	manager2 "github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/lobby/config"
	"github.com/jzyong/ugk/lobby/handler"
	"github.com/jzyong/ugk/lobby/manager"
	"github.com/jzyong/ugk/lobby/rpc"
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

	handler.RegisterClientHandler() // 没引用handler不执行init，手动执行一下

	m.Run()
	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()
}

type ModuleManager struct {
	*util.DefaultModuleManager
	LobbyManager         *manager.LobbyManager
	GateKcpClientManager *manager2.GateKcpClientManager
	PlayerManager        *manager.PlayerManager
	GrpcManager          *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.LobbyManager = m.AppendModule(manager.GetLobbyManager()).(*manager.LobbyManager)
	m.GateKcpClientManager = m.AppendModule(manager2.GetGateKcpClientManager()).(*manager2.GateKcpClientManager)
	m.PlayerManager = m.AppendModule(manager.GetPlayerManager()).(*manager.PlayerManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
