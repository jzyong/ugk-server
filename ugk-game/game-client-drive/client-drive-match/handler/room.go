package handler

//// 进入房间
//func enterRoom(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
//	request := &message.GalacticKittensEnterRoomRequest{}
//	err := proto.Unmarshal(msg.Bytes, request)
//	response := &message.GalacticKittensEnterRoomResponse{}
//	if err != nil {
//		log.Error("解析消息错误：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
//		return
//	}
//
//	for i, p := range room.Players {
//		if p.Id == msg.PlayerId {
//			//移除之前的，从新进入，可能断网重进，网关等连接已经改变了等
//			log.Info("玩家：%v已进入房间", msg.PlayerId)
//			room.Players = append(room.Players[:i], room.Players[i+1:]...)
//		}
//	}
//
//	// 需要向大厅获取玩家基础信息 ,暂时只考虑只有一个lobby，后面修改
//	hallGrpc, err, lobbyId := manager.GetServiceClientManager().GetLobbyGrpcByPlayerId(msg.PlayerId)
//	if err != nil {
//		log.Error("获取大厅异常：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
//		return
//	}
//	client := message.NewPlayerServiceClient(hallGrpc)
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
//	defer cancel()
//	playerInfoResponse, err := client.GetPlayerInfo(ctx, &message.PlayerInfoRequest{PlayerId: msg.PlayerId})
//	if err != nil {
//		log.Error("请求玩家信息：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
//		return
//	}
//
//	player = mode.NewPlayer(msg.PlayerId)
//	player.GateClient = gateClient
//	player.SetHeartTime(time.Now())
//	playerInfo := playerInfoResponse.GetPlayer()
//	player.Level = playerInfo.GetLevel()
//	player.Exp = playerInfo.GetExp()
//	player.Nick = playerInfo.GetNick()
//	player.LobbyId = lobbyId
//	player.CharacterId = int32(len(room.Players))
//	room.Players = append(room.Players, player)
//	response.Result = &message.MessageResult{
//		Status: 200,
//		Msg:    "success",
//	}
//
//	gateClient.SendToGate(player.Id, message.MID_GalacticKittensEnterRoomRes, response, msg.Seq)
//	manager2.GetRoomManager().BroadcastRoomInfo(room)
//	log.Debug("%v进入房间%v", player.Id, room.Id)
//}
//
//// 退出房间
//func quitRoom(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
//	response := &message.GalacticKittensQuitRoomResponse{}
//	if player == nil {
//		log.Error("%v未登录", msg.PlayerId)
//		response.Result = &message.MessageResult{
//			Status: 404,
//			Msg:    "Player not login",
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensQuitRoomRes, response, msg.Seq)
//		return
//	}
//	request := &message.GalacticKittensQuitRoomRequest{}
//	err := proto.Unmarshal(msg.Bytes, request)
//	if err != nil {
//		log.Error("解析消息错误：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensQuitRoomRes, response, msg.Seq)
//		return
//	}
//
//	for i, p := range room.Players {
//		if p.Id == player.Id {
//			room.Players = append(room.Players[:i], room.Players[i+1:]...)
//			break
//		}
//	}
//	response.Result = &message.MessageResult{
//		Status: 200,
//		Msg:    "success",
//	}
//
//	gateClient.SendToGate(player.Id, message.MID_GalacticKittensQuitRoomRes, response, msg.Seq)
//	manager2.GetRoomManager().BroadcastRoomInfo(room)
//	log.Debug("%v退出房间%v", player.Id, room.Id)
//}
//
//// 选择角色
//func selectCharacter(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
//	response := &message.GalacticKittenSelectCharacterResponse{}
//	if player == nil {
//		log.Error("%v未登录", msg.PlayerId)
//		response.Result = &message.MessageResult{
//			Status: 404,
//			Msg:    "Player not login",
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittenSelectCharacterRes, response, msg.Seq)
//		return
//	}
//	request := &message.GalacticKittenSelectCharacterRequest{}
//	err := proto.Unmarshal(msg.Bytes, request)
//	if err != nil {
//		log.Error("解析消息错误：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittenSelectCharacterRes, response, msg.Seq)
//		return
//	}
//
//	player.CharacterId = request.GetCharacterId()
//	response.Result = &message.MessageResult{
//		Status: 200,
//		Msg:    "success",
//	}
//
//	gateClient.SendToGate(player.Id, message.MID_GalacticKittenSelectCharacterRes, response, msg.Seq)
//	manager2.GetRoomManager().BroadcastRoomInfo(room)
//	log.Debug("%v选择角色%v", player.Id, player.CharacterId)
//}
//
//// 准备
//func prepare(player *mode.Player, room *mode.Room, gateClient *manager.GateKcpClient, msg *mode2.UgkMessage) {
//	response := &message.GalacticKittensPrepareResponse{}
//	if player == nil {
//		log.Error("%v未登录", msg.PlayerId)
//		response.Result = &message.MessageResult{
//			Status: 404,
//			Msg:    "Player not login",
//		}
//		gateClient.SendToGate(msg.PlayerId, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
//		return
//	}
//	request := &message.GalacticKittensPrepareRequest{}
//	err := proto.Unmarshal(msg.Bytes, request)
//
//	if err != nil {
//		log.Error("解析消息错误：%v", err)
//		response.Result = &message.MessageResult{
//			Status: 500,
//			Msg:    err.Error(),
//		}
//		gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
//		return
//	}
//
//	//设置准备状态
//	if request.Prepare && room.StateMachine.IsInState(manager2.InitStateRoom) {
//		room.StateMachine.ChangeState(manager2.PrepareStateRoom)
//	}
//	player.Prepare = request.Prepare
//	gateClient.SendToGate(player.Id, message.MID_GalacticKittensPrepareRes, response, msg.Seq)
//	// 推送房间消息
//	manager2.GetRoomManager().BroadcastRoomInfo(room)
//
//	//检测是否可进入游戏
//	prepareCount := 0
//	for _, p := range room.Players {
//		if p.Prepare {
//			prepareCount++
//		}
//	}
//	if prepareCount == 0 { //玩家退出房间这些暂时不考虑
//		room.StateMachine.ChangeState(manager2.InitStateRoom)
//	} else if prepareCount == len(room.Players) {
//		room.StateMachine.ChangeState(manager2.LoadStateRoom)
//	}
//}
