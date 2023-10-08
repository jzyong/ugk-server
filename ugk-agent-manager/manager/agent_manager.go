package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
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
}

func (m *AgentManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
