package manager

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/lobby/config"
	"github.com/jzyong/ugk/lobby/mode"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// 消息执行函数
type handFunc func(player *mode.Player, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage)

// GateHandlers 客户端消息处理器
var GateHandlers = make(map[uint32]handFunc)

// LobbyManager  入口
type LobbyManager struct {
	util.DefaultModule
}

var lobbyManager *LobbyManager
var lobbySingletonOnce sync.Once

func GetLobbyManager() *LobbyManager {
	lobbySingletonOnce.Do(func() {
		lobbyManager = &LobbyManager{}
	})
	return lobbyManager
}

func (m *LobbyManager) Init() error {
	log.Info("LobbyManager 初始化......")

	manager.GetZookeeperManager().Start(config.BaseConfig)
	m.writeZooKeeperConfigs()
	manager.GetGateKcpClientManager().Start(config.BaseConfig)
	manager.GetGateKcpClientManager().ServerHeartRequest = &message.ServerHeartRequest{Server: &message.ServerInfo{
		Id:   config.BaseConfig.Id,
		Name: config.BaseConfig.Name,
	}}
	return nil
}

// 向zookeeper写配置数据，暂时全在这里写入
func (m *LobbyManager) writeZooKeeperConfigs() {
	// excel 配置
	mongoExcelConfig := &config2.MongoConfig{
		Url:      config.BaseConfig.MongoUrl,
		Database: "ugk-config",
	}
	util.ZKUpdate(manager.GetZookeeperManager().GetConn(), fmt.Sprintf(config2.ZKMongoExcelConfigPath, config.BaseConfig.Profile), util.ToString(mongoExcelConfig))

	//redis数据库
	util.ZKUpdate(manager.GetZookeeperManager().GetConn(), fmt.Sprintf(config2.ZKRedisPath, config.BaseConfig.Profile), util.ToString(config.BaseConfig.Redis))
}

func (m *LobbyManager) Run() {
}

func (m *LobbyManager) Stop() {
	manager.GetZookeeperManager().Stop()
}
