package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ServerService 服务器rpc请求
type ServerService struct {
	message.UnimplementedServerServiceServer
}

func (service *ServerService) CloseServer(ctx context.Context, request *message.CloseServerRequest) (*message.CloseServerResponse, error) {

	response := &message.CloseServerResponse{}
	log.Info("服务器关闭")
	//TODO 添加逻辑，注册方法执行？

	return response, nil
}

func (service *ServerService) ReloadConfig(ctx context.Context, request *message.ReloadConfigRequest) (*message.ReloadConfigResponse, error) {

	//TODO 添加逻辑，注册方法执行？
	return nil, status.Errorf(codes.Unimplemented, "method ReloadConfig not implemented")
}
