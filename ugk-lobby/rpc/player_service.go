package rpc

import (
	"context"
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
	response := &message.PlayerInfoResponse{Player: info}
	return response, nil
}
