package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/gate/manager"
	"github.com/jzyong/ugk/message/message"
)

// 心跳
func heart(user *manager.User, data []byte, seq uint32, timeStamp int64) {
	log.Info("%d 心跳 序号=%d 时间=%d", user.Id, seq, timeStamp)
	// 返回消息
	user.SendToClient(message.MID_HeartRes, &message.HeartResponse{}, seq)
}

// 登录
func login(user *manager.User, data []byte, seq uint32, timeStamp int64) {
	log.Info("%d 登录 序号=%d 时间=%d", user.Id, seq, timeStamp)
}
