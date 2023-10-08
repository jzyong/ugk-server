package rpc

import (
	"fmt"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/golib/util"
	"github.com/jzyong/ugk/agent-manager/config"
	"github.com/jzyong/ugk/common/rpc"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/grpc"
	"net"
	"strings"
)

// GRpcManager grpc 管理
type GRpcManager struct {
	*util.DefaultModule
	GrpcServer *grpc.Server
}

func (m *GRpcManager) Init() error {
	server := grpc.NewServer()
	m.GrpcServer = server
	// 添加grpc服务
	message.RegisterServerServiceServer(server, new(rpc.ServerService))
	//容器中运行 绑定ip地址不一定正确走配置
	portStr := strings.Split(config.BaseConfig.RpcUrl, ":")[1]
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", portStr))
	if err != nil {
		log.Fatal("%v", err)
	}
	log.Info("grpc listen on:%v", config.BaseConfig.RpcUrl)
	go server.Serve(listen)
	log.Info("GRpcManager:init end......")

	return nil
}

func (m *GRpcManager) Run() {

}

func (m *GRpcManager) Stop() {
	if m.GrpcServer != nil {
		m.GrpcServer.Stop()
	}
}
