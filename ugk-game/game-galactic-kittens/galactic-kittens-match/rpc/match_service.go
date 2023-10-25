package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/galactic-kittens-match/config"
	manager2 "github.com/jzyong/ugk/galactic-kittens-match/manager"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// MatchService 匹配rpc请求
type MatchService struct {
	message.UnimplementedGalacticKittensMatchServiceServer
}

// PlayerServerList 玩家服务器列表 TODO 待测试
func (service *ServerService) PlayerServerList(ctx context.Context, request *message.GalacticKittensPlayerServerListRequest) (*message.GalacticKittensPlayerServerListResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	response := &message.GalacticKittensPlayerServerListResponse{}

	manager2.GetRoomManager().ProcessFun <- func() {

		room := manager2.GetRoomManager().GetRoom(request.GetRoomId())
		gateServers := make(map[int64]*message.ServerInfo, len(room.Players))
		lobbyServers := make(map[int64]*message.ServerInfo, len(room.Players))
		room.ProcessFun <- func() {
			defer wg.Done()

			lobbyClients := manager.GetServiceClientManager().GetClients(config2.GetZKServicePath(config.BaseConfig.Profile, config2.LobbyName, 0))
			for _, player := range room.Players {
				gateServer := &message.ServerInfo{
					Id:      player.GateClient.Id,
					Name:    config2.GateName,
					GrpcUrl: player.GateClient.Url,
				}
				gateServers[player.Id] = gateServer

				// 大厅 暂时只有一个，先这样 TODO
				if lobbyClients != nil {
					for _, c := range lobbyClients {
						serverInfo := &message.ServerInfo{
							Id:      c.Id,
							Name:    config2.LobbyName,
							GrpcUrl: c.Url,
						}
						lobbyServers[player.Id] = serverInfo
						break
					}
				}
			}
			response.PlayerGateServers = gateServers
			response.PlayerLobbyServers = lobbyServers
		}
		wg.Done()
	}

	wg.Wait()

	log.Info("服务信息：%v", response)
	return response, nil
}

// GameFinish 游戏完成 TODO 待测试
func (service *ServerService) GameFinish(ctx context.Context, request *message.GalacticKittensGameFinishRequest) (*message.GalacticKittensGameFinishResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	response := &message.GalacticKittensGameFinishResponse{}

	manager2.GetRoomManager().ProcessFun <- func() {
		room := manager2.GetRoomManager().GetRoom(request.GetRoomId())
		room.ProcessFun <- func() {
			defer wg.Done()
			close(room.GetCloseChan())
			response.Result = &message.MessageResult{
				Status: 200,
				Msg:    "success",
			}
		}
		wg.Done()
	}

	wg.Wait()

	log.Info("服务信息：%v", response)
	return response, nil
	return nil, nil
}
