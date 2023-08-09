package handler

import (
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
)

//func init() {
//	registerClientHandler()
//}

func RegisterClientHandler() {
	//登录模块
	manager.ClientHandlers[uint32(message.MID_HeartReq)] = heart
	manager.ClientHandlers[uint32(message.MID_LoginReq)] = login

}
