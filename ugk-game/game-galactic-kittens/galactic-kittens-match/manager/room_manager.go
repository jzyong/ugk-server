package manager

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/golib/util/fsm"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"sync"
	"time"
)

// RoomManager 房间  每个房间一个routine //TODO  待测试
type RoomManager struct {
	util.DefaultModule
	IdRooms     map[uint32]*mode.Room  //房间 key：房间id
	PlayerRooms map[int64]uint32       //房间 key：玩家ID
	messages    chan *mode2.UgkMessage //收到的所有消息
}

var roomManager *RoomManager
var roomSingletonOnce sync.Once

func GetRoomManager() *RoomManager {
	roomSingletonOnce.Do(func() {
		roomManager = &RoomManager{
			IdRooms:     make(map[uint32]*mode.Room, 64),
			PlayerRooms: make(map[int64]uint32, 1024),
			messages:    make(chan *mode2.UgkMessage, 1024),
		}
	})
	return roomManager
}

func (m *RoomManager) Init() error {
	log.Info("RoomManager 初始化......")
	//设置消息处理
	manager.GetGateKcpClientManager().MessageHandFunc = m.messageDistribute
	return nil
}
func (m *RoomManager) Run() {
	go m.run()
}

func (m *RoomManager) Stop() {
}

func (m *RoomManager) run() {
	for {
		select {
		case message := <-m.messages: //转发消息到房间
			room := m.GetRoomByPlayerId(message.PlayerId)
			room.GetMessages() <- message
		}
	}
}

func (m *RoomManager) GetRoomByPlayerId(playerId int64) *mode.Room {
	if roomId, ok := m.PlayerRooms[playerId]; ok {
		return m.GetRoom(roomId)
	} else {
		// 自动分配房间
		return m.autoAssignRoom(playerId)
	}
}

// 自动分配房间
func (m *RoomManager) autoAssignRoom(playerId int64) *mode.Room {
	for _, room := range m.IdRooms {
		if len(room.Players) < 4 && room.StateMachine.IsInState(roomInitState) {
			return room
		}
	}
	server := GetDataManager().GetServer()
	server.RoomId += 1
	return m.GetRoom(server.RoomId)
}

func (m *RoomManager) GetRoom(id uint32) *mode.Room {
	if room, ok := m.IdRooms[id]; ok {
		return room
	} else {
		room = mode.NewRoom(id)
		room.StateMachine = &fsm.DefaultStateMachine[*mode.Room]{Owner: room}
		room.StateMachine.SetInitialState(roomInitState)
		m.IdRooms[id] = room
		go roomRun(room)
		return room
	}
}

// 消息分发处理
func (m *RoomManager) messageDistribute(playerId int64, msg *mode2.UgkMessage) {
	// 转发到房间管理器routine执行
	m.messages <- msg

}

// 运行玩家routine
func roomRun(room *mode.Room) {
	secondTicker := time.Tick(time.Second)
	for {
		select {
		case msg := <-room.GetMessages(): //消息处理
			handRequest(room, msg)
		case <-secondTicker:
			roomSecondUpdate(room)
		case <-room.GetCloseChan():
			log.Info("房间：%d 关闭", room.Id)
			return
		}

	}
}

func handRequest(room *mode.Room, msg *mode2.UgkMessage) {
	defer mode2.ReturnUgkMessage(msg)
	handFunc := GateHandlers[msg.MessageId]
	if handFunc == nil {
		log.Warn("消息：%d未实现，玩家%d逻辑处理失败", msg.MessageId, room.Id)
		return
	}
	var player *mode.Player
	for _, p := range room.Players {
		if p.Id == msg.PlayerId {
			player = p
		}
	}

	handFunc(player, room, msg.Client.(*manager.GateKcpClient), msg)
	room.SetHeartTime(util.Now())
	log.Debug("%d 收到消息 mid=%d seq=%d", room.Id, msg.MessageId, msg.Seq)
}

// 房间每秒监测
func roomSecondUpdate(room *mode.Room) {
	if util.Now().Sub(room.GetHeartTime()) > config2.ServerHeartInterval {
		log.Info("房间：%d 心跳超时离线：%v", room.Id, util.Now().Sub(room.GetHeartTime()).Minutes())
		close(room.GetCloseChan())
	}
	room.StateMachine.Update()
}
