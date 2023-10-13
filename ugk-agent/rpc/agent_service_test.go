package rpc

import (
	"context"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

const grpcUrl = "192.168.110.2:3031"

// 创建游戏容器
func TestCreateGameContainer(t *testing.T) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(grpcUrl, dialOption)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	client := message.NewAgentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := &message.CreateGameServiceRequest{
		GameId:         1,
		GameName:       "game-galactic-kittens",
		ControlGrpcUrl: "192.168.110.2:4000",
	}

	response, err := client.CreateGameService(ctx, request)

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("resonse:%v", response)
}

// 关闭游戏容器
func TestCloseGameContainer(t *testing.T) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(grpcUrl, dialOption)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	client := message.NewAgentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := &message.CloseGameServiceRequest{
		GameId:   1,
		GameName: "game-galactic-kittens",
	}

	response, err := client.CloseGameService(ctx, request)

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("resonse:%v", response)
}
