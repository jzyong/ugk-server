package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"sync"
)

// PlayerManager 玩家
type PlayerManager struct {
	util.DefaultModule
}

var playerManager *PlayerManager
var playerSingletonOnce sync.Once

func GetPlayerManager() *PlayerManager {
	playerSingletonOnce.Do(func() {
		playerManager = &PlayerManager{}
	})
	return playerManager
}

func (m *PlayerManager) Init() error {
	log.Info("PlayerManager 初始化......")
	return nil
}
func (m *PlayerManager) Run() {
}

func (m *PlayerManager) Stop() {
}
