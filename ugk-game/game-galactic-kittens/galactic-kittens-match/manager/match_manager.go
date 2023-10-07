package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// 消息执行函数
type handFunc func(player *mode.Player, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage)

// GateHandlers 客户端消息处理器
var GateHandlers = make(map[uint32]handFunc)

// MatchManager  入口
type MatchManager struct {
	util.DefaultModule
}

var matchSingletonOnce sync.Once
var matchManager *MatchManager

func GetMatchManager() *MatchManager {
	matchSingletonOnce.Do(func() {
		matchManager = &MatchManager{}
	})
	return matchManager
}

func (m *MatchManager) Init() error {
	log.Info("MatchManager 初始化......")
	manager.GetZookeeperManager().Start(config.BaseConfig)
	manager.GetGateKcpClientManager().Start(config.BaseConfig)
	manager.GetGateKcpClientManager().ServerHeartRequest = &message.ServerHeartRequest{Server: &message.ServerInfo{
		Id:   config.BaseConfig.Id,
		Name: config.BaseConfig.Name,
	}}
	return nil
}

func (m *MatchManager) Run() {
}

func (m *MatchManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
