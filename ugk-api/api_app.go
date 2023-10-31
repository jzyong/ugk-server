package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/api/config"
	"github.com/jzyong/ugk/api/controller"
	"github.com/jzyong/ugk/api/manager"
	"github.com/jzyong/ugk/api/rpc"
	"runtime"
)

// 登录入口
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfigAndLog()
	log.Info("启动 api ......")

	controller.RegisterController()

	var err error
	err = m.Init()
	if err != nil {
		log.Error("api 启动错误: %s", err.Error())
		return
	}

	m.Run()
	util.WaitForTerminate()
	m.Stop()

}

type ModuleManager struct {
	*util.DefaultModuleManager
	ApiManager  *manager.ApiManager
	WebManager  *manager.WebManager
	GrpcManager *rpc.GRpcManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.ApiManager = m.AppendModule(manager.GetApiManager()).(*manager.ApiManager)
	m.WebManager = m.AppendModule(manager.GetWebManager()).(*manager.WebManager)
	m.GrpcManager = m.AppendModule(&rpc.GRpcManager{}).(*rpc.GRpcManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
