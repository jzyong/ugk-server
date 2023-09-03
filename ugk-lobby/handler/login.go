package handler

import (
	"github.com/jzyong/golib/log"
	"github.com/jzyong/ugk/lobby/mode"
	"github.com/jzyong/ugk/message/message"
	"google.golang.org/protobuf/proto"
)

// 加载玩家数据
func loadPlayer(player *mode.Player, data []byte, seq uint32) {
	request := &message.LoadPlayerRequest{}
	err := proto.Unmarshal(data, request)
	if err != nil {
		log.Error("解析消息错误：%v", err)
		//TODO 返给客户端错误信息
		return
	}
	log.Info("%d 加载玩家数据", request.GetPlayerId())
	//TODO 封装消息返回
}
