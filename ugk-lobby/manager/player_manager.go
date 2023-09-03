package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/lobby/mode"
	"sync"
)

// PlayerManager 玩家 TODO 每个玩家一个routine
type PlayerManager struct {
	util.DefaultModule
	IdPlayers map[int64]*mode.Player //登录后的玩家ID用户
}

var playerManager *PlayerManager
var playerSingletonOnce sync.Once

func GetPlayerManager() *PlayerManager {
	playerSingletonOnce.Do(func() {
		playerManager = &PlayerManager{
			IdPlayers: make(map[int64]*mode.Player, 1024),
		}
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
