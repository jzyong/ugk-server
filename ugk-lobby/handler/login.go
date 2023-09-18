package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/common/manager"
	mode2 "github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/lobby/mode"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
)

// 加载玩家数据
func loadPlayer(player *mode.Player, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
	request := &message.LoadPlayerRequest{}
	err := proto.Unmarshal(msg.Bytes, request)
	response := &message.LoadPlayerResponse{}
	if err != nil {
		log.Error("解析消息错误：%v", err)
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		gateClient.SendToGate(player.Id, message.MID_LoadPlayerRes, response, msg.Seq)
		return
	}
	log.Info("%d 加载玩家数据", request.GetPlayerId())
	//TODO 其他数据
	response.PlayerInfo = &message.PlayerInfo{
		PlayerId: player.Id,
		Nick:     player.Nick,
		Level:    player.Level,
		Exp:      player.Exp,
		Gold:     10000000,
	}

	gateClient.SendToGate(player.Id, message.MID_LoadPlayerRes, response, msg.Seq)

}