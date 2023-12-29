package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/client-drive-match/config"
	config2 "github.com/jzyong/ugk/common/config"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/common/rpc"
	"github.com/jzyong/ugk/message/message"
)

// ServerService 服务器rpc请求
type ServerService struct {
	rpc.ServerService
}

// GetServerInfo 获取网关（kcp）和大厅服务信息，子游戏前期调试需要
func (service *ServerService) GetServerInfo(ctx context.Context, request *message.GetServerInfoRequest) (*message.GetServerInfoResponse, error) {
	response := &message.GetServerInfoResponse{}
	serverInfos := make([]*message.ServerInfo, 0, 4)

	//网关地址 暂时连接所有网关，后面只连接有玩家的网关
	gateClients := manager.GetGateKcpClientManager().Clients
	for _, c := range gateClients {
		serverInfo := &message.ServerInfo{
			Id:      c.Id,
			Name:    config2.GateName,
			GrpcUrl: c.Url,
		}
		serverInfos = append(serverInfos, serverInfo)
	}

	//大厅
	lobbyClients := manager.GetServiceClientManager().GetClients(config2.GetZKServicePath(config.BaseConfig.Profile, config2.LobbyName, 0))
	if lobbyClients != nil {
		for _, c := range lobbyClients {
			serverInfo := &message.ServerInfo{
				Id:      c.Id,
				Name:    config2.LobbyName,
				GrpcUrl: c.Url,
			}
			serverInfos = append(serverInfos, serverInfo)
		}
	} else {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    "lobby service not find",
		}
	}
	if len(serverInfos) < 1 {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    "lobby and gate service not find",
		}
	} else {
		response.ServerInfo = serverInfos
	}

	log.Info("服务信息：%v", response)
	return response, nil
}
