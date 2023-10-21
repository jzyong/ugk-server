package mode

import (
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/mode"
	"github.com/xtaci/kcp-go/v5"
	"time"
)

// Player 玩家，每个玩家一个routine处理逻辑 TODO 玩家发送消息，routine消息处理，记得关闭
type Player struct {
	Id          int64                 `_id`   //唯一id
	Nick        string                `nick`  //昵称
	Level       uint32                `level` //等级
	Exp         uint32                `exp`   //经验
	Gold        uint64                `gold`  // 金币
	Items       map[uint32]Item       `items` //道具
	gateSession *kcp.UDPSession       //网关连接会话
	messages    chan *mode.UgkMessage //接收到的玩家消息
	closeChan   chan struct{}         //离线等关闭Chan
	heartTime   time.Time             //心跳时间
	dirty       bool                  //数据是否改变
}

func NewPlayer(id int64) *Player {
	player := &Player{
		Id:        id,
		messages:  make(chan *mode.UgkMessage, 1024),
		closeChan: make(chan struct{}),
		heartTime: util.Now(),
	}
	return player
}

func (player *Player) GetGateSession() *kcp.UDPSession {
	return player.gateSession
}

func (player *Player) SetGateSession(session *kcp.UDPSession) {
	player.gateSession = session
}

// GetMessages 待处理的玩家消息
func (player *Player) GetMessages() chan *mode.UgkMessage {
	return player.messages
}

func (player *Player) SetMessages(messages chan *mode.UgkMessage) {
	player.messages = messages
}

func (player *Player) GetCloseChan() chan struct{} {
	return player.closeChan
}

func (player *Player) SetCloseChan(closeChan chan struct{}) {
	player.closeChan = closeChan
}

func (player *Player) GetHeartTime() time.Time {
	return player.heartTime
}
func (player *Player) SetHeartTime(time time.Time) {
	player.heartTime = time
}

func (player *Player) SetDirty(dirty bool) {
	player.dirty = true
}

func (player *Player) GetDirty() bool {
	return player.dirty
}
