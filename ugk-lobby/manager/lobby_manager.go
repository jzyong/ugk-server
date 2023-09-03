package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"sync"
)

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
	return nil
}

func (m *LobbyManager) Run() {
}

func (m *LobbyManager) Stop() {
}
