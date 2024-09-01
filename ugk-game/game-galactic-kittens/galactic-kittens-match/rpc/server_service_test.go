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

// 获取服务器信息
func TestGetServerInfo(t *testing.T) {
	dialOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial("127.0.0.1:4000", dialOption)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	client := message.NewServerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	response, err := client.GetServerInfo(ctx, &message.GetServerInfoRequest{})

	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("resonse:%v", response)
	time.Sleep(time.Second)
}
