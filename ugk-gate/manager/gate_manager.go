package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/gate/config"
	"sync"
)

// GateManager  入口
type GateManager struct {
	util.DefaultModule
}

var gateManager *GateManager
var gateSingletonOnce sync.Once

func GetGateManager() *GateManager {
	gateSingletonOnce.Do(func() {
		gateManager = &GateManager{}
	})
	return gateManager
}

func (m *GateManager) Init() error {
	log.Info("GateManager 初始化......")
	manager.GetZookeeperManager().Start(config.BaseConfig)
	return nil
}

func (m *GateManager) Run() {
}

func (m *GateManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
