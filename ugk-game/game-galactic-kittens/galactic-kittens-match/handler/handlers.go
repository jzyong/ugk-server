package handler

import (
	"github.com/jzyong/ugk/galactic-kittens-match/manager"
	"github.com/jzyong/ugk/message/message"
)

// RegisterClientHandler 注册客户端处理消息
func RegisterClientHandler() {
	manager.GateHandlers[uint32(message.MID_GalacticKittensEnterRoomReq)] = enterRoom
	manager.GateHandlers[uint32(message.MID_GalacticKittensQuitRoomReq)] = quitRoom
	manager.GateHandlers[uint32(message.MID_GalacticKittenSelectCharacterReq)] = selectCharacter
	manager.GateHandlers[uint32(message.MID_GalacticKittensPrepareReq)] = prepare

}
