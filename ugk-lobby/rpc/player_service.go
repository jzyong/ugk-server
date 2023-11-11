package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/lobby/manager"
	"github.com/jzyong/ugk/message/message"
)

// PlayerService 玩家 grpc
type PlayerService struct {
	message.UnimplementedPlayerServiceServer
}

func (m *PlayerService) GetPlayerInfo(ctx context.Context, request *message.PlayerInfoRequest) (*message.PlayerInfoResponse, error) {
	player := manager.GetPlayerManager().GetPlayer(request.GetPlayerId())

	info := &message.PlayerInfo{
		PlayerId: player.Id,
		Nick:     player.Nick,
		Level:    player.Level,
		Exp:      player.Exp,
		Gold:     player.Gold,
	}
	log.Debug("请求%v 角色信息：%v", request.GetPlayerId(), info)
	response := &message.PlayerInfoResponse{Player: info, Result: &message.MessageResult{
		Status: 200,
		Msg:    "success",
	}}
	return response, nil
}
