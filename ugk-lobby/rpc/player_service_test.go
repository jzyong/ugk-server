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

const grpcUrl = "192.168.110.2:3021"

// 获取玩家信息
func TestGetPlayerInfo(t *testing.T) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(grpcUrl, dialOption)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	client := message.NewPlayerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	request := &message.PlayerInfoRequest{PlayerId: 1}

	response, err := client.GetPlayerInfo(ctx, request)

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("resonse:%v", response)
	time.Sleep(time.Second * 10)
}
