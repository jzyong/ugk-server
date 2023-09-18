package rpc

import (
	"context"
	"github.com/jzyong/ugk/login/manager"
	"github.com/jzyong/ugk/message/message"
)

// LoginService 聊天rpc请求
type LoginService struct {
	message.UnimplementedLoginServiceServer
}

var accounts = make(map[string]int64)

func (service *LoginService) Login(ctx context.Context, request *message.LoginRequest) (*message.LoginResponse, error) {
	response := &message.LoginResponse{}

	account := manager.GetDataManager().FindAccount(request.GetAccount())

	if account == nil {
		response.Result = &message.MessageResult{
			Status: 404,
			Msg:    "Account password error",
		}
	} else {
		response.Result = &message.MessageResult{
			Status: 200,
			Msg:    "Login Success",
		}
		response.PlayerId = account.PlayerId
	}

	return response, nil
}
