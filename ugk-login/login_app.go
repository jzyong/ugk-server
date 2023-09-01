package main

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/login/config"
	"github.com/jzyong/ugk/login/manager"
	"runtime"
)

// 登录入口
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	config.InitConfigAndLog()
	log.Info("启动 login ......")

	var err error
	err = m.Init()
	if err != nil {
		log.Error("login 启动错误: %s", err.Error())
		return
	}

	m.Run()
	util.WaitForTerminate()
	m.Stop()

	util.WaitForTerminate()
}

type ModuleManager struct {
	*util.DefaultModuleManager
	LoginManager *manager.LoginManager
}

// Init 初始化模块
func (m *ModuleManager) Init() error {
	m.LoginManager = m.AppendModule(manager.GetLoginManager()).(*manager.LoginManager)
	return m.DefaultModuleManager.Init()
}

var m = &ModuleManager{
	DefaultModuleManager: util.NewDefaultModuleManager(),
}
