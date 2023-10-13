package rpc

import (
	"context"
	"github.com/jzyong/ugk/agent/manager"
	"github.com/jzyong/ugk/message/message"
)

// AgentService agent
type AgentService struct {
	message.UnimplementedAgentServiceServer
}

func (a *AgentService) CreateGameService(ctx context.Context, request *message.CreateGameServiceRequest) (*message.CreateGameServiceResponse, error) {
	response := manager.GetDockerManager().CreateGameServiceContainer(request)
	return response, nil
}

func (a *AgentService) CloseGameService(ctx context.Context, request *message.CloseGameServiceRequest) (*message.CloseGameServiceResponse, error) {
	response := manager.GetDockerManager().CloseGameServiceContainer(request)
	return response, nil
}
