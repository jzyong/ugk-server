package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/message/message"
)

// AgentControlService agent controller
type AgentControlService struct {
	message.UnimplementedAgentControlServiceServer
}

func (a *AgentControlService) HostMachineInfoUpload(ctx context.Context, request *message.HostMachineInfoUploadRequest) (*message.HostMachineInfoUploadResponse, error) {

	//TODO 需要同步，放在一个routine中执行
	log.Info("请求主机信息：%v", request)

	return nil, nil
}
