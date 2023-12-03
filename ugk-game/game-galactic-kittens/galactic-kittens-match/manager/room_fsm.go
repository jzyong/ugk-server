package manager

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/golib/util/fsm"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"github.com/jzyong/ugk/message/message"
	"time"
)

//房间 状态机
//初始化，准备，加载，游戏中，完成，结束

var InitStateRoom = &RoomInitState{}
var PrepareStateRoom = &RoomPrepareState{}
var LoadStateRoom = &RoomLoadState{}
var GamingStateRoom = &RoomGamingState{}
var FinishStateRoom = &RoomFinishState{}
var CloseStateRoom = &RoomCloseState{}

// RoomInitState 初始化
type RoomInitState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomPrepareState 准备
type RoomPrepareState struct {
	fsm.EmptyState[*mode.Room]
}

// RoomLoadState 加载
type RoomLoadState struct {
	fsm.EmptyState[*mode.Room]
}

func (state *RoomLoadState) Enter(room *mode.Room) {
	GetRoomManager().BroadcastRoomInfo(room)

	//编辑器运行模式直接返回，不创建docker容器
	if config.EditorMode {
		room.StateMachine.ChangeState(GamingStateRoom)
		return
	}

	// 请求agent-manager创建游戏服务
	grpcClient, err := manager.GetServiceClientManager().GetGrpc(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentManagerName, 0), 0)
	if err != nil {
		//TODO 需要做异常处理
		log.Error("%v创建房间游戏服失败：%v", room.Id, err)
		return
	}
	client := message.NewAgentControlServiceClient(grpcClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	request := &message.CreateGameServiceRequest{
		GameId:         room.Id,
		GameName:       "game-galactic-kittens",
		ControlGrpcUrl: config.BaseConfig.GetRpcUrl(),
	}
	response, err := client.CreateGameService(ctx, request)
	if err != nil {
		//TODO 需要做异常处理
		log.Error("%v创建房间游戏服失败：%v", room.Id, err)
		return
	}
	if response.GetResult().GetStatus() == 200 {
		room.StateMachine.ChangeState(GamingStateRoom)
	} else {
		//TODO 需要做异常处理
		log.Error("%v创建房间游戏服失败：%v", room.Id, response.GetResult().GetMsg())
	}

}

// RoomGamingState 游戏中
type RoomGamingState struct {
	fsm.EmptyState[*mode.Room]
}

func (state *RoomGamingState) Enter(room *mode.Room) {
	GetRoomManager().BroadcastRoomInfo(room)
}

// RoomFinishState 完成
type RoomFinishState struct {
	fsm.EmptyState[*mode.Room]
}

func (state *RoomFinishState) Enter(room *mode.Room) {
	GetRoomManager().BroadcastRoomInfo(room)
	room.CloseTime = util.Now().Add(time.Minute)
	//TODO 向大厅同步更新玩家游戏结果
}

func (state *RoomFinishState) Update(room *mode.Room) {
	// 倒计时（60s）变更为关闭状态，TODO 或者所有玩家退出房间
	if util.Now().After(room.CloseTime) {
		room.StateMachine.ChangeState(CloseStateRoom)
	}
}

// RoomCloseState 关闭
type RoomCloseState struct {
	fsm.EmptyState[*mode.Room]
}

func (state *RoomCloseState) Enter(room *mode.Room) {
	//客户端收到该消息，将所有玩家踢回大厅
	GetRoomManager().BroadcastRoomInfo(room)
	// 请求agent-manager创建游戏服务
	grpcClient, err := manager.GetServiceClientManager().GetGrpc(config2.GetZKServicePath(config.BaseConfig.Profile, config2.AgentManagerName, 0), 0)
	if err != nil {
		//TODO 需要做异常处理
		log.Error("%v关闭房间游戏服失败：%v", room.Id, err)
		return
	}
	client := message.NewAgentControlServiceClient(grpcClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()
	request := &message.CloseGameServiceRequest{
		GameId:   room.Id,
		GameName: "game-galactic-kittens",
	}
	response, err := client.CloseGameService(ctx, request)
	if err != nil {
		//TODO 需要做异常处理
		log.Error("%v关闭房间游戏服失败：%v", room.Id, err)
		return
	}
	if response.GetResult().GetStatus() == 200 {
		log.Info("%v关闭房间游戏服成功", room.Id)
	} else {
		//TODO 需要做异常处理
		log.Error("%v关闭房间游戏服失败：%v", room.Id, response.GetResult().GetMsg())
	}

	// 关闭房间
	close(room.GetCloseChan())

}

// RoomState 房间状态
func RoomState(room *mode.Room) uint32 {
	if room.StateMachine.IsInState(InitStateRoom) {
		return 0
	} else if room.StateMachine.IsInState(PrepareStateRoom) {
		return 1
	} else if room.StateMachine.IsInState(LoadStateRoom) {
		return 2
	} else if room.StateMachine.IsInState(GamingStateRoom) {
		return 3
	} else if room.StateMachine.IsInState(FinishStateRoom) {
		return 4
	} else if room.StateMachine.IsInState(CloseStateRoom) {
		return 5
	}
	return 0

}
