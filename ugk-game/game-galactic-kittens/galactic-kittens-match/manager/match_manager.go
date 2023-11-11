package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// 消息执行函数
type handFunc func(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage)

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
	//写消息ID
	messageIds := make([]uint32, 0, len(GateHandlers))
	for messageId, _ := range GateHandlers {
		messageIds = append(messageIds, messageId)
	}
	messageIdPath := fmt.Sprintf(config2.ZKMessageIdPath, config.BaseConfig.Profile, config2.GameGalacticKittensMatch)
	log.Debug("注册消息：%v", messageIds)
	util.ZKUpdate(manager.GetZookeeperManager().GetConn(), messageIdPath, util.ToString(messageIds))
	return nil
}

func (m *MatchManager) Run() {
	//监听并连接微服务
	manager.GetServiceClientManager().WatchGrpcService(config2.GetZKServicePath(config.BaseConfig.Profile, config2.LobbyName, 0))
	manager.GetServiceClientManager().WatchGrpcService(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentManagerName, 0))
}

func (m *MatchManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
