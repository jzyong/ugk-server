package mode

import (
	"github.com/xtaci/kcp-go/v5"
)

// Player 玩家，每个玩家一个routine处理逻辑 TODO 玩家发送消息，routine消息处理，记得关闭
type Player struct {
	Id          int64           `id`    //唯一id
	Nick        string          `nick`  //昵称
	Level       uint32          `level` //等级
	Exp         uint32          `exp`   //经验
	Items       map[uint32]Item `items` //道具
	gateSession *kcp.UDPSession //网关连接会话
}

func NewPlayer(id int64) *Player {
	player := &Player{Id: id}
	return player
}

func (player *Player) GetGateSession() *kcp.UDPSession {
	return player.gateSession
}

func (player *Player) SetGateSession(session *kcp.UDPSession) {
	player.gateSession = session
}
