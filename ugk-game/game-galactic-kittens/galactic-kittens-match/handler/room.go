package handler

import (
	"context"
	"github.com/jzyong/golib/log"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	manager2 "github.com/jzyong/ugk/galactic-kittens-match/manager"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
	"time"
)

// 进入房间 TODO 待测试
func enterRoom(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
	request := &message.GalacticKittensEnterRoomRequest{}
	err := proto.Unmarshal(msg.Bytes, request)
	response := &message.GalacticKittensEnterRoomResponse{}
	if err != nil {
		log.Error("解析消息错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(request.GetPlayerId(), message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
		return
	}

	// 需要向大厅获取玩家基础信息 ,暂时只考虑只有一个lobby，后面修改
	hallGrpc, err := manager.GetServiceClientManager().GetGrpcConn(config2.GetZKServicePath(config.BaseConfig.Profile, config2.LobbyName, 0), 0)
	if err != nil {
		log.Error("获取大厅异常：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(request.GetPlayerId(), message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
		return
	}
	client := message.NewPlayerServiceClient(hallGrpc)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	playerInfoResponse, err := client.GetPlayerInfo(ctx, &message.PlayerInfoRequest{PlayerId: request.GetPlayerId()})
	if err != nil {
		log.Error("请求玩家信息：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(request.GetPlayerId(), message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
		return
	}

	player = mode.NewPlayer(request.GetPlayerId())
	player.GateClient = gateClient
	player.SetHeartTime(time.Now())
	playerInfo := playerInfoResponse.GetPlayer()
	player.Level = playerInfo.GetLevel()
	player.Exp = playerInfo.GetExp()
	player.Nick = playerInfo.GetNick()
	room.Players = append(room.Players, player)
	response.Result = &message.MessageResult{
		Status: 200,
		Msg:    "success",
	}

	gateClient.SendToGate(player.Id, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
	manager2.GetRoomManager().BroadcastRoomInfo(room)
}

// 准备 TODO 待测试
func prepare(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
	response := &message.GalacticKittensPrepareResponse{}
	if player == nil {
		log.Error("%未登录：%v", msg.PlayerId)
		response.Result = &message.MessageResult{
			Status: 404,
			Msg:    "Player not login",
		}
		gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
		return
	}
	request := &message.GalacticKittensPrepareRequest{}
	err := proto.Unmarshal(msg.Bytes, request)

	if err != nil {
		log.Error("解析消息错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
		return
	}

	//设置准备状态
	if request.Prepare && room.StateMachine.IsInState(manager2.InitStateRoom) {
		room.StateMachine.ChangeState(manager2.PrepareStateRoom)
	}
	player.Prepare = request.Prepare

	prepareCount := 0
	for _, p := range room.Players {
		if p.Prepare {
			prepareCount++
		}
	}
	if prepareCount == 0 { //玩家退出房间这些暂时不考虑
		room.StateMachine.ChangeState(manager2.InitStateRoom)
	} else if prepareCount == len(room.Players) {
		room.StateMachine.ChangeState(manager2.LoadStateRoom)
	}

	gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
	// 推送房间消息
	manager2.GetRoomManager().BroadcastRoomInfo(room)
}
