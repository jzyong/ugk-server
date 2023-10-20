package mode

import (
	"github.com/jzyong/golib/util/fsm"
	"github.com/jzyong/ugk/common/mode"
	"time"
)

// Room 房间 同一房间在同一个routine中执行
type Room struct {
	Id           uint32                  //房间ID
	Players      []*Player               //玩家
	StateMachine fsm.StateMachine[*Room] //状态机
	messages     chan *mode.UgkMessage   //接收到的玩家消息
	closeChan    chan struct{}           //离线等关闭Chan
	heartTime    time.Time               //心跳时间
}

func NewRoom(id uint32) *Room {
	room := &Room{
		Id:        id,
		messages:  make(chan *mode.UgkMessage, 1024),
		closeChan: make(chan struct{}),
	}
	return room
}

// GetMessages 待处理的玩家消息
func (room *Room) GetMessages() chan *mode.UgkMessage {
	return room.messages
}

func (room *Room) GetCloseChan() chan struct{} {
	return room.closeChan
}

func (room *Room) GetHeartTime() time.Time {
	return room.heartTime
}
func (room *Room) SetHeartTime(time time.Time) {
	room.heartTime = time
}
