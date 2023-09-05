package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/lobby/config"
	"github.com/jzyong/ugk/lobby/mode"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// 消息执行函数
type handFunc func(player *mode.Player, data []byte, seq uint32)

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

	manager.GetGateKcpClientManager().ConnectKcpServer(config.AppConfigManager.GateUrl)
	manager.GetGateKcpClientManager().ServerHeartRequest = &message.ServerHeartRequest{Server: &message.ServerInfo{
		Id:   config.AppConfigManager.Id,
		Name: "lobby",
	}}
	return nil
}

func (m *LobbyManager) Run() {
}

func (m *LobbyManager) Stop() {
}
