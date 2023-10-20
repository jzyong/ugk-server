package mode

import (
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
	"time"
)

// Player 玩家，每个玩家一个routine处理逻辑 TODO 玩家发送消息，routine消息处理，记得关闭 ;提取公共的到ugk-common中
type Player struct {
	Id         int64                  `id`    //唯一id
	Nick       string                 `nick`  //昵称
	Level      uint32                 `level` //等级
	Exp        uint32                 `exp`   //经验
	Prepare    bool                   //是否准备
	GateClient *manager.GateKcpClient //网关客户端
	heartTime  time.Time              //心跳时间
}

func NewPlayer(id int64) *Player {
	player := &Player{
		Id:        id,
		heartTime: util.Now(),
	}
	return player
}

func (player *Player) GetHeartTime() time.Time {
	return player.heartTime
}
func (player *Player) SetHeartTime(time time.Time) {
	player.heartTime = time
}

// SendMsg 主推消息
func (player *Player) SendMsg(mid message.MID, msg proto.Message) {
	player.GateClient.SendToGate(player.Id, mid, msg, 0)
}
