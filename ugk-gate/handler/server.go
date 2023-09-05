package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
	"time"
)

// 服务器心跳
func serverHeart(user *manager.User, data []byte, seq uint32, client *manager.GameKcpClient) {
	client.HeartTime = time.Now()
	request := &message.ServerHeartRequest{}
	response := &message.ServerHeartResponse{}
	err := proto.Unmarshal(data, request)
	if err != nil {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		log.Error("消息返序列化错误：%v", err)
		//TODO 发送返回消息
		client.SendToGame(0, message.MID_ServerHeartRes, response, 0)
		return
	}
	serverInfo := request.GetServer()
	log.Info("%v-%v 心跳", serverInfo.GetId(), serverInfo.GetName())
	client.SendToGame(0, message.MID_ServerHeartRes, response, 0)
}
