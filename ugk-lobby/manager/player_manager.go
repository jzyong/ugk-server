package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/lobby/mode"
	"sync"
	"time"
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
	manager.GetGateKcpClientManager().MessageHandFunc = m.messageDistribute
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
		go playerRun(player)
		return player
	}
}

// 消息分发处理
func (m *PlayerManager) messageDistribute(playerId int64, msg *mode2.UgkMessage) {
	player := m.GetPlayer(playerId)
	player.GetMessages() <- msg

}

// 运行玩家routine
func playerRun(player *mode.Player) {
	//TODO 记得关闭
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case msg := <-player.GetMessages(): //消息处理
			handRequest(player, msg)
		case <-secondTicker:
			playerSecondUpdate(player)
		case <-player.GetCloseChan():
			//TODO 存储玩家数据等离线操作
			log.Info("玩家：%d 离线", player.Id)
			return
		}

	}
}

func handRequest(player *mode.Player, msg *mode2.UgkMessage) {
	defer mode2.ReturnUgkMessage(msg)
	handFunc := GateHandlers[msg.MessageId]
	if handFunc == nil {
		log.Warn("消息：%d未实现，玩家%d逻辑处理失败", msg.MessageId, player.Id)
		return
	}
	handFunc(player, msg.Client.(*manager.GateKcpClient), msg)
	player.SetHeartTime(util.Now())
	log.Debug("%d 收到消息 mid=%d seq=%d", player.Id, msg.MessageId, msg.Seq)
}

// 玩家每秒监测
func playerSecondUpdate(player *mode.Player) {
	if util.Now().Sub(player.GetHeartTime()) > config2.ServerHeartInterval {
		log.Info("玩家：%d 心跳超时离线：%v", player.Id, util.Now().Sub(player.GetHeartTime()).Minutes())
		close(player.GetCloseChan())
	}
}
