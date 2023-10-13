package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/agent-manager/manager"
	"github.com/jzyong/ugk/message/message"
	"sync"
)

// AgentControlService agent controller
type AgentControlService struct {
	message.UnimplementedAgentControlServiceServer
}

func (a *AgentControlService) HostMachineInfoUpload(ctx context.Context, request *message.HostMachineInfoUploadRequest) (*message.HostMachineInfoUploadResponse, error) {
	log.Trace("请求主机信息：%v", request)
	manager.GetDockerManager().MachiInfoChan <- request.GetHostMachineInfo()

	return &message.HostMachineInfoUploadResponse{Result: &message.MessageResult{
		Status: 200,
		Msg:    "成功",
	}}, nil
}

// TODO 待测试
func (a *AgentControlService) CreateGameService(ctx context.Context, request *message.CreateGameServiceRequest) (*message.CreateGameServiceResponse, error) {
	var wg sync.WaitGroup
	wg.Add(2)
	var response = &message.CreateGameServiceResponse{}
	manager.GetDockerManager().RequestChan <- func() {
		manager.GetDockerManager().CreateGameService(ctx, wg, request, response)
	}
	return response, nil
}

func (a *AgentControlService) CloseGameService(ctx context.Context, request *message.CloseGameServiceRequest) (*message.CloseGameServiceResponse, error) {
	//TODO
	return nil, nil
}
