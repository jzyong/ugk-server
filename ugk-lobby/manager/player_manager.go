package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/lobby/mode"
	"sync"
)

// PlayerManager 玩家 TODO 每个玩家一个routine
type PlayerManager struct {
	util.DefaultModule
	IdPlayers map[int64]*mode.Player //登录后的玩家ID用户
	mutex     sync.RWMutex           //玩家锁
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
	//设置消息处理
	manager.GetGateKcpClientManager().MessageHandFunc = m.messageHand
	return nil
}
func (m *PlayerManager) Run() {
}

func (m *PlayerManager) Stop() {
}

func (m *PlayerManager) GetPlayer(id int64) *mode.Player {
	defer m.mutex.RUnlock()
	m.mutex.RLock()
	if player, ok := m.IdPlayers[id]; ok {
		return player
	} else {
		//TODO 从数据库查询
		player = mode.NewPlayer(id)
		m.IdPlayers[id] = player
		return player
	}
}

// 消息分发处理
func (m *PlayerManager) messageHand(playerId int64, messageId uint32, seq uint32, timeStamp int64, data []byte, client *manager.GateKcpClient) {
	//player := m.GetPlayer(playerId)
	//TODO 转发到玩家routine
	log.Info("%d 收到消息 mid=%d seq=%d", playerId, messageId, seq)
}
