package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	manager2 "github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	"github.com/jzyong/ugk/galactic-kittens-match/handler"
	"github.com/jzyong/ugk/galactic-kittens-match/manager"
	"github.com/jzyong/ugk/galactic-kittens-match/rpc"
	"runtime"
)

// gate 入口
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfigAndLog()
	log.Info("启动 galactic kittens match ......")

	handler.RegisterClientHandler() // 没引用handler不执行init，手动执行一下
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
	GateManager          *manager.MatchManager
	ServiceClientManager *manager2.ServiceClientManager
	GrpcManager          *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.GateManager = m.AppendModule(manager.GetMatchManager()).(*manager.MatchManager)
	m.ServiceClientManager = m.AppendModule(manager2.GetServiceClientManager()).(*manager2.ServiceClientManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
