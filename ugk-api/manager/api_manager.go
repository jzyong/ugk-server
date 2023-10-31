package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/api/config"
	"github.com/jzyong/ugk/common/manager"
	"sync"
)

// ApiManager  入口
type ApiManager struct {
	util.DefaultModule
}

var apiManager *ApiManager
var apiSingletonOnce sync.Once

func GetApiManager() *ApiManager {
	apiSingletonOnce.Do(func() {
		apiManager = &ApiManager{}
	})
	return apiManager
}

func (m *ApiManager) Init() error {
	log.Info("ApiManager 初始化......")
	manager.GetZookeeperManager().Start(config.BaseConfig)
	return nil
}

func (m *ApiManager) Run() {
}

func (m *ApiManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
