package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/gate/config"
	"github.com/jzyong/ugk/gate/manager"
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
	m.Run()
	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()
}

type ModuleManager struct {
	*util.DefaultModuleManager
	GateManager   *manager.GateManager
	ClientManager *manager.ClientManager
	UserManager   *manager.UserManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.GateManager = m.AppendModule(manager.GetGateManager()).(*manager.GateManager)
	m.ClientManager = m.AppendModule(manager.GetClientManager()).(*manager.ClientManager)
	m.UserManager = m.AppendModule(manager.GetUserManager()).(*manager.UserManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
