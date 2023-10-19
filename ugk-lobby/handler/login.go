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
	response.PlayerInfo = &message.PlayerInfo{
		PlayerId: player.Id,
		Nick:     player.Nick,
		Level:    player.Level,
		Exp:      player.Exp,
		Gold:     10000000,
	}

	games := make([]*message.GameInfo, 0, 3)
	//TODO 走配置
	game1 := &message.GameInfo{
		GameId: 1,
		Name:   "GalacticKittens",
		Status: 0,
		Icon:   "",
	}
	games = append(games, game1)
	game2 := &message.GameInfo{
		GameId: 2,
		Name:   "Coon",
		Status: 0,
		Icon:   "",
	}
	games = append(games, game2)
	game3 := &message.GameInfo{
		GameId: 3,
		Name:   "Racing",
		Status: 0,
		Icon:   "",
	}
	games = append(games, game3)
	response.GameInfo = games

	gateClient.SendToGate(player.Id, message.MID_LoadPlayerRes, response, msg.Seq)

}
