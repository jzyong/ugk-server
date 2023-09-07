package handler

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/common/mode"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
)

// 心跳
func heart(user *manager.User, msg *mode.UgkMessage) {
	log.Info("%d 心跳 序号=%d 时间=%d", user.Id, msg.Seq, msg.TimeStamp)
	// 返回消息
	user.SendToClient(message.MID_HeartRes, &message.HeartResponse{}, msg.Seq)
}

// 登录
func login(user *manager.User, msg *mode.UgkMessage) {
	request := &message.LoginRequest{}
	proto.Unmarshal(msg.Bytes, request)

	//TODO 通过zookeeper获取login服务
	log.Info("%d 登录 序号=%d %+v", user.Id, msg.Seq, request)
	conn := manager.GetLoginClientManager().ClientConn
	client := message.NewLoginServiceClient(conn)
	response, err := client.Login(context.Background(), request)
	if err != nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}}, msg.Seq)
		return
	}
	user.Id = response.PlayerId

	lobbyClient := manager.GetServerManager().AssignLobby(user.Id)
	if lobbyClient == nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 404,
			Msg:    "lobby server not open",
		}}, msg.Seq)
		return
	}
	user.LobbyClient = lobbyClient
	manager.GetUserManager().AddUser(user)

	user.SendToClient(message.MID_LoginRes, response, msg.Seq)

}
