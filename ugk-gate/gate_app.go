package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	manager2 "github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/gate/config"
	"github.com/jzyong/ugk/gate/handler"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/gate/rpc"
	"runtime"
)

// gate 入口
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
	handler.RegisterGameHandler()

	m.Run()
	util.WaitForTerminate()
	m.Stop()

}

type ModuleManager struct {
	*util.DefaultModuleManager
	GateManager          *manager.GateManager
	ClientManager        *manager.ClientManager
	ServerManager        *manager.ServerManager
	UserManager          *manager.UserManager
	ServiceClientManager *manager2.ServiceClientManager
	LoginClientManager   *manager.LoginClientManager
	GrpcManager          *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.GateManager = m.AppendModule(manager.GetGateManager()).(*manager.GateManager)
	m.ClientManager = m.AppendModule(manager.GetClientManager()).(*manager.ClientManager)
	m.ServerManager = m.AppendModule(manager.GetServerManager()).(*manager.ServerManager)
	m.UserManager = m.AppendModule(manager.GetUserManager()).(*manager.UserManager)
	m.ServiceClientManager = m.AppendModule(manager2.GetServiceClientManager()).(*manager2.ServiceClientManager)
	m.LoginClientManager = m.AppendModule(manager.GetLoginClientManager()).(*manager.LoginClientManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
