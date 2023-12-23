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

// PlayerServerList 玩家服务器列表
func (service *MatchService) PlayerServerList(ctx context.Context, request *message.GalacticKittensPlayerServerListRequest) (*message.GalacticKittensPlayerServerListResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	response := &message.GalacticKittensPlayerServerListResponse{}

	manager2.GetRoomManager().ProcessFun <- func() {
		roomId := request.GetRoomId()
		if roomId < 1 {
			roomId = manager2.GetDataManager().GetServer().RoomId
		}

		room := manager2.GetRoomManager().GetRoom(roomId)
		gateServers := make(map[int64]*message.ServerInfo, len(room.Players))
		lobbyServers := make(map[int64]*message.ServerInfo, len(room.Players))
		playerInfos := make(map[int64]*message.GalacticKittensPlayerServerListResponse_PlayerInfo, len(room.Players))
		room.ProcessFun <- func() {
			defer wg.Done()

			if len(room.Players) < 1 {
				response.Result = &message.MessageResult{
					Status: 500,
					Msg:    "room no player",
				}
				return
			}

			for _, player := range room.Players {
				gateServer := &message.ServerInfo{
					Id:      player.GateClient.Id,
					Name:    config2.GateName,
					GrpcUrl: player.GateClient.Url,
				}
				gateServers[player.Id] = gateServer
				_, err, lobbyId := manager.GetServiceClientManager().GetLobbyGrpcByPlayerId(player.Id)
				if err != nil {
					log.Warn("%v 没有正确获得大厅", player.Id)
					continue
				}

				lobbyClient := manager.GetServiceClientManager().GetClient(config2.GetZKServicePath(config.BaseConfig.Profile, config2.LobbyName, 0), lobbyId)
				if lobbyClient != nil {
					serverInfo := &message.ServerInfo{
						Id:      lobbyClient.Id,
						Name:    config2.LobbyName,
						GrpcUrl: lobbyClient.Url,
					}
					lobbyServers[player.Id] = serverInfo
				}

				playerInfo := &message.GalacticKittensPlayerServerListResponse_PlayerInfo{CharacterId: player.CharacterId}
				playerInfos[player.Id] = playerInfo

			}
			response.PlayerGateServers = gateServers
			response.PlayerLobbyServers = lobbyServers
			response.PlayerInfos = playerInfos
			response.RoomId = roomId
		}
		wg.Done()
	}

	wg.Wait()

	log.Info("服务信息：%v", response)
	return response, nil
}

// GameFinish 游戏完成
func (service *MatchService) GameFinish(ctx context.Context, request *message.GalacticKittensGameFinishRequest) (*message.GalacticKittensGameFinishResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	response := &message.GalacticKittensGameFinishResponse{}
	log.Info("房间 %d 结束:%v", request.GetRoomId(), request)
	manager2.GetRoomManager().ProcessFun <- func() {
		room := manager2.GetRoomManager().GetRoom(request.GetRoomId())
		room.ProcessFun <- func() {
			defer wg.Done()

			// 分数计算，推送给客户端,传回大厅等
			clientResponse := &message.GalacticKittensGameFinishResponse{}
			clientResponse.Victory = request.Victory

			bestIndex := 0
			var bestScore uint32 = 0
			statistics := make([]*message.GalacticKittensGameFinishResponse_PlayerStatistics, 0, len(request.Statistics))
			for i, statistic := range request.Statistics {
				info := &message.GalacticKittensGameFinishResponse_PlayerStatistics{
					PlayerId:      statistic.GetPlayerId(),
					KillCount:     statistic.GetKillCount(),
					Score:         statistic.GetKillCount()*100 - statistic.GetUsePowerCount()*50,
					UsePowerCount: statistic.GetUsePowerCount(),
					Best:          false,
					Victory:       statistic.GetVictory(),
				}
				if info.KillCount > bestScore {
					bestScore = info.KillCount
					bestIndex = i
				}
				statistics = append(statistics, info)
			}
			statistics[bestIndex].Best = true
			clientResponse.Statistics = statistics
			manager2.GetRoomManager().BroadcastMsg(room, message.MID_GalacticKittensGameFinishRes, clientResponse)
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
