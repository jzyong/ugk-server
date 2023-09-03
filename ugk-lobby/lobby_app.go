package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/lobby/config"
	"github.com/jzyong/ugk/lobby/handler"
	"github.com/jzyong/ugk/lobby/manager"
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
	LobbyManager  *manager.LobbyManager
	NetManager    *manager.NetManager
	PlayerManager *manager.PlayerManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.LobbyManager = m.AppendModule(manager.GetLobbyManager()).(*manager.LobbyManager)
	m.NetManager = m.AppendModule(manager.GetNetManager()).(*manager.NetManager)
	m.PlayerManager = m.AppendModule(manager.GetPlayerManager()).(*manager.PlayerManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
