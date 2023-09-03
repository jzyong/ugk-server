package rpc

import (
	"context"
	"github.com/jzyong/ugk/message/message"
)

// LoginService 聊天rpc请求
type LoginService struct {
	message.UnimplementedLoginServiceServer
}

var accounts = make(map[string]int64)

func init() {
	accounts["test1"] = 1
	accounts["test2"] = 2
	accounts["test3"] = 3
	accounts["test4"] = 4

}

func (service *LoginService) Login(ctx context.Context, request *message.LoginRequest) (*message.LoginResponse, error) {
	//TODO 暂时写死4个账号
	response := &message.LoginResponse{}

	if id, ok := accounts[request.Account]; !ok {
		response.Result = &message.MessageResult{
			Status: 404,
			Msg:    "Account password error",
		}
	} else {
		response.Result = &message.MessageResult{
			Status: 200,
			Msg:    "Login Success",
		}
		response.PlayerId = id
	}

	return response, nil
}
