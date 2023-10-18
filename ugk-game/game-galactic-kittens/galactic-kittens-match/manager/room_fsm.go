package manager

import (
	"github.com/jzyong/golib/util/fsm"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
)

//房间 状态机
//初始化，准备，加载，游戏中，完成，结束

var roomInitState = &RoomInitState{}
var roomPrepareState = &RoomPrepareState{}
var roomLoadState = &RoomLoadState{}
var roomGamingState = &RoomGamingState{}
var roomFinishState = &RoomFinishState{}
var roomCloseState = &RoomCloseState{}

// RoomInitState 初始化
type RoomInitState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomPrepareState 准备
type RoomPrepareState struct {
	fsm.EmptyState[*mode.Room]
}

func (r *RoomPrepareState) Update(room *mode.Room) {

}

// RoomLoadState 加载
type RoomLoadState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomGamingState 游戏中
type RoomGamingState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomFinishState 完成
type RoomFinishState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomCloseState 关闭
type RoomCloseState struct {
	fsm.EmptyState[*mode.Room]
}
