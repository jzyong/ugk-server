package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/agent-manager/manager"
	"github.com/jzyong/ugk/message/message"
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
