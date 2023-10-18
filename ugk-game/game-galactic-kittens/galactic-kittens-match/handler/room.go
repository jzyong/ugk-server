package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/galactic-kittens-match/mode"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
)

// 进入房间
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
		gateClient.SendToGate(player.Id, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
		return
	}

	//TODO

	gateClient.SendToGate(player.Id, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)

}

// 准备
func prepare(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
	request := &message.GalacticKittensPrepareRequest{}
	err := proto.Unmarshal(msg.Bytes, request)
	response := &message.GalacticKittensPrepareResponse{}
	if err != nil {
		log.Error("解析消息错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
		return
	}

	//TODO 

	gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)

}
