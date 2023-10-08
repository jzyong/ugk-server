package rpc

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/common/manager"
	"github.com/jzyong/ugk/common/rpc"
	"github.com/jzyong/ugk/message/message"
)

// ServerService 服务器rpc请求
type ServerService struct {
	rpc.ServerService
}

// GetServerInfo 获取网关（kcp）和大厅服务信息，子游戏需要
func (service *ServerService) GetServerInfo(ctx context.Context, request *message.GetServerInfoRequest) (*message.GetServerInfoResponse, error) {

	response := &message.GetServerInfoResponse{}
	manager.GetZookeeperManager()

	log.Info("服务器关闭")
	//TODO 添加逻辑，注册方法执行？

	return response, nil
}
