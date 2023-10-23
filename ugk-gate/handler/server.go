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
		client.SendToGame(0, message.MID_ServerHeartRes, response, msg.Seq)
		return
	}
	serverInfo := request.GetServer()
	client.Id = serverInfo.GetId()
	client.Name = serverInfo.GetName()
	manager.GetServerManager().UpdateGameServer(serverInfo, client)
	log.Trace("%v-%v 心跳 时间=%v", serverInfo.GetId(), serverInfo.GetName(), msg.TimeStamp)
	client.SendToGame(0, message.MID_ServerHeartRes, response, msg.Seq)
}

// 子游戏通知网关绑定玩家网络连接 TODO 待测试
func bindGameConnect(user *manager.User, client *manager.GameKcpClient, msg *mode.UgkMessage) {
	request := &message.BindGameConnectRequest{}
	response := &message.BindGameConnectResponse{}
	err := proto.Unmarshal(msg.Bytes, request)
	if err != nil {
		response.Result = &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}
		log.Error("消息返序列化错误：%v", err)
		client.SendToGame(user.Id, message.MID_BindGameConnectRes, response, 0)
		return
	}
	if request.GetBind() {
		user.GameClient = client
		log.Debug("玩家：%d 进入游戏%v", user.Id, client.Name)
	} else {
		user.GameMessages = nil
		log.Debug("玩家：%d 退出游戏%v", user.Id, client.Name)
	}

	response.Result = &message.MessageResult{Status: 200, Msg: "成功"}
	client.SendToGame(user.Id, message.MID_BindGameConnectRes, response, 0)
}
