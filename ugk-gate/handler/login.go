package handler

import (
	"context"
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
)

// 心跳
func heart(user *manager.User, data []byte, seq uint32, timeStamp int64) {
	log.Info("%d 心跳 序号=%d 时间=%d", user.Id, seq, timeStamp)
	// 返回消息
	user.SendToClient(message.MID_HeartRes, &message.HeartResponse{}, seq)
}

// 登录
func login(user *manager.User, data []byte, seq uint32, timeStamp int64) {
	request := &message.LoginRequest{}
	proto.Unmarshal(data, request)

	//TODO 通过zookeeper获取login服务
	log.Info("%d 登录 序号=%d %+v", user.Id, seq, request)
	conn := manager.GetLoginClientManager().ClientConn
	client := message.NewLoginServiceClient(conn)
	response, err := client.Login(context.Background(), request)
	if err != nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}}, seq)
		return
	}
	user.Id = response.PlayerId

	lobbyClient := manager.GetServerManager().AssignLobby(user.Id)
	if lobbyClient == nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 404,
			Msg:    "lobby server not open",
		}}, seq)
		return
	}
	user.LobbyClient = lobbyClient
	manager.GetUserManager().AddUser(user)

	user.SendToClient(message.MID_LoginRes, response, seq)

}
