package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
	"time"
)

// 服务器心跳
func serverHeart(user *manager.User, client *manager.GameKcpClient, msg *mode.UgkMessage) {
	client.HeartTime = time.Now()
	request := &message.ServerHeartRequest{}
	response := &message.ServerHeartResponse{}
	err := proto.Unmarshal(msg.Bytes, request)
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
	client.Id = serverInfo.GetId()
	client.Name = serverInfo.GetName()
	manager.GetServerManager().UpdateGameServer(serverInfo, client)
	log.Trace("%v-%v 心跳 时间=%v", serverInfo.GetId(), serverInfo.GetName(), msg.TimeStamp)
	client.SendToGame(0, message.MID_ServerHeartRes, response, 0)
}
