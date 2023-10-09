package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"sync"
)

// AgentManager  入口
type AgentManager struct {
	util.DefaultModule
}

var agentManager *AgentManager
var agentSingletonOnce sync.Once

func GetAgentManager() *AgentManager {
	agentSingletonOnce.Do(func() {
		agentManager = &AgentManager{}
	})
	return agentManager
}

func (m *AgentManager) Init() error {
	log.Info("AgentManager 初始化......")

	manager.GetZookeeperManager().Start(config.BaseConfig)
	return nil
}

func (m *AgentManager) Run() {

	//监听并连接agent服务
	manager.GetServiceClientManager().WatchGrpcService(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentName, 0))
}

func (m *AgentManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
