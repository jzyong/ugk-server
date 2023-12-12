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
	//log.Trace("%d 心跳 序号=%d 时间=%d", user.Id, msg.Seq, msg.TimeStamp)

	request := &message.HeartRequest{}
	proto.Unmarshal(msg.Bytes, request)

	// 如果在unity游戏中，需要转发unity服务器  TODO 待测试
	if user.GameClient != nil {
		user.GameClient.SendToGame(user.Id, message.MID_HeartReq, request, msg.Seq)
		return
	}
	// 返回消息
	user.SendToClient(message.MID_HeartRes, &message.HeartResponse{ClientTime: request.GetClientTime()}, msg.Seq)
}

// 登录
func login(user *manager.User, msg *mode.UgkMessage) {
	request := &message.LoginRequest{}
	proto.Unmarshal(msg.Bytes, request)

	log.Info("%d 登录 序号=%d %+v", user.Id, msg.Seq, request)
	// 不能随机获取登录服，存在客户端多次请求发送创建账号，可能同一设备到两个登录服去创建账号了
	loginClient := manager.GetLoginClientManager().HashModClient(user.ClientSession.RemoteAddr().String())
	if loginClient == nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 500,
			Msg:    "login service close",
		}}, msg.Seq)
		return
	}

	client := message.NewLoginServiceClient(loginClient.ClientConn)
	response, err := client.Login(context.Background(), request)
	if err != nil {
		user.SendToClient(message.MID_LoginRes, &message.LoginResponse{Result: &message.MessageResult{
			Status: 500,
			Msg:    err.Error(),
		}}, msg.Seq)
		log.Error("%v 登录错误:%v", request.GetAccount(), err)
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
