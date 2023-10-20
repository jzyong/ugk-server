// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.18.1
// source: game_galactic_kittens.proto

package message

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 进入房间
type GalacticKittensEnterRoomRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId int64 `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"` //玩家id
}

func (x *GalacticKittensEnterRoomRequest) Reset() {
	*x = GalacticKittensEnterRoomRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensEnterRoomRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensEnterRoomRequest) ProtoMessage() {}

func (x *GalacticKittensEnterRoomRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensEnterRoomRequest.ProtoReflect.Descriptor instead.
func (*GalacticKittensEnterRoomRequest) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{0}
}

func (x *GalacticKittensEnterRoomRequest) GetPlayerId() int64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

// 进入房间
type GalacticKittensEnterRoomResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *MessageResult `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` //错误信息
}

func (x *GalacticKittensEnterRoomResponse) Reset() {
	*x = GalacticKittensEnterRoomResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensEnterRoomResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensEnterRoomResponse) ProtoMessage() {}

func (x *GalacticKittensEnterRoomResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensEnterRoomResponse.ProtoReflect.Descriptor instead.
func (*GalacticKittensEnterRoomResponse) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{1}
}

func (x *GalacticKittensEnterRoomResponse) GetResult() *MessageResult {
	if x != nil {
		return x.Result
	}
	return nil
}

// 推送房间信息
type GalacticKittensRoomInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Room *GalacticKittensRoomInfo `protobuf:"bytes,1,opt,name=room,proto3" json:"room,omitempty"` //房间信息
}

func (x *GalacticKittensRoomInfoResponse) Reset() {
	*x = GalacticKittensRoomInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensRoomInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensRoomInfoResponse) ProtoMessage() {}

func (x *GalacticKittensRoomInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensRoomInfoResponse.ProtoReflect.Descriptor instead.
func (*GalacticKittensRoomInfoResponse) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{2}
}

func (x *GalacticKittensRoomInfoResponse) GetRoom() *GalacticKittensRoomInfo {
	if x != nil {
		return x.Room
	}
	return nil
}

// 房间信息
type GalacticKittensRoomInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64                        `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` //房间ID
	Player []*GalacticKittensPlayerInfo `protobuf:"bytes,2,rep,name=player,proto3" json:"player,omitempty"`
	State  int32                        `protobuf:"varint,3,opt,name=state,proto3" json:"state,omitempty"` //房间状态 0匹配；1准备完成；2游戏中；3结算；4游戏结束
}

func (x *GalacticKittensRoomInfo) Reset() {
	*x = GalacticKittensRoomInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensRoomInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensRoomInfo) ProtoMessage() {}

func (x *GalacticKittensRoomInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensRoomInfo.ProtoReflect.Descriptor instead.
func (*GalacticKittensRoomInfo) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{3}
}

func (x *GalacticKittensRoomInfo) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GalacticKittensRoomInfo) GetPlayer() []*GalacticKittensPlayerInfo {
	if x != nil {
		return x.Player
	}
	return nil
}

func (x *GalacticKittensRoomInfo) GetState() int32 {
	if x != nil {
		return x.State
	}
	return 0
}

// 玩家信息
type GalacticKittensPlayerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId     int64  `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"`         //玩家id
	Nick         string `protobuf:"bytes,2,opt,name=nick,proto3" json:"nick,omitempty"`                  //昵称
	Prepare      bool   `protobuf:"varint,3,opt,name=prepare,proto3" json:"prepare,omitempty"`           //是否准备
	Score        int32  `protobuf:"varint,4,opt,name=score,proto3" json:"score,omitempty"`               //分数
	PowerUpCount int32  `protobuf:"varint,5,opt,name=powerUpCount,proto3" json:"powerUpCount,omitempty"` //充能数
	Hp           int32  `protobuf:"varint,6,opt,name=hp,proto3" json:"hp,omitempty"`                     //血量
	Icon         string `protobuf:"bytes,7,opt,name=icon,proto3" json:"icon,omitempty"`                  //头像
}

func (x *GalacticKittensPlayerInfo) Reset() {
	*x = GalacticKittensPlayerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensPlayerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensPlayerInfo) ProtoMessage() {}

func (x *GalacticKittensPlayerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensPlayerInfo.ProtoReflect.Descriptor instead.
func (*GalacticKittensPlayerInfo) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{4}
}

func (x *GalacticKittensPlayerInfo) GetPlayerId() int64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *GalacticKittensPlayerInfo) GetNick() string {
	if x != nil {
		return x.Nick
	}
	return ""
}

func (x *GalacticKittensPlayerInfo) GetPrepare() bool {
	if x != nil {
		return x.Prepare
	}
	return false
}

func (x *GalacticKittensPlayerInfo) GetScore() int32 {
	if x != nil {
		return x.Score
	}
	return 0
}

func (x *GalacticKittensPlayerInfo) GetPowerUpCount() int32 {
	if x != nil {
		return x.PowerUpCount
	}
	return 0
}

func (x *GalacticKittensPlayerInfo) GetHp() int32 {
	if x != nil {
		return x.Hp
	}
	return 0
}

func (x *GalacticKittensPlayerInfo) GetIcon() string {
	if x != nil {
		return x.Icon
	}
	return ""
}

// 准备
type GalacticKittensPrepareRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Prepare bool `protobuf:"varint,1,opt,name=prepare,proto3" json:"prepare,omitempty"` //ture准备，false取消
}

func (x *GalacticKittensPrepareRequest) Reset() {
	*x = GalacticKittensPrepareRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensPrepareRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensPrepareRequest) ProtoMessage() {}

func (x *GalacticKittensPrepareRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensPrepareRequest.ProtoReflect.Descriptor instead.
func (*GalacticKittensPrepareRequest) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{5}
}

func (x *GalacticKittensPrepareRequest) GetPrepare() bool {
	if x != nil {
		return x.Prepare
	}
	return false
}

// 准备
type GalacticKittensPrepareResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *MessageResult `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"` //错误信息
}

func (x *GalacticKittensPrepareResponse) Reset() {
	*x = GalacticKittensPrepareResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensPrepareResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensPrepareResponse) ProtoMessage() {}

func (x *GalacticKittensPrepareResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensPrepareResponse.ProtoReflect.Descriptor instead.
func (*GalacticKittensPrepareResponse) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{6}
}

func (x *GalacticKittensPrepareResponse) GetResult() *MessageResult {
	if x != nil {
		return x.Result
	}
	return nil
}

// 进入游戏 内部
type GalacticKittensEnterGameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId int64 `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"` //玩家id
}

func (x *GalacticKittensEnterGameRequest) Reset() {
	*x = GalacticKittensEnterGameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensEnterGameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensEnterGameRequest) ProtoMessage() {}

func (x *GalacticKittensEnterGameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensEnterGameRequest.ProtoReflect.Descriptor instead.
func (*GalacticKittensEnterGameRequest) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{7}
}

func (x *GalacticKittensEnterGameRequest) GetPlayerId() int64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

// 进入游戏 内部
type GalacticKittensEnterGameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *MessageResult             `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` //错误信息
	Player *GalacticKittensPlayerInfo `protobuf:"bytes,2,opt,name=player,proto3" json:"player,omitempty"` //玩家
}

func (x *GalacticKittensEnterGameResponse) Reset() {
	*x = GalacticKittensEnterGameResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensEnterGameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensEnterGameResponse) ProtoMessage() {}

func (x *GalacticKittensEnterGameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensEnterGameResponse.ProtoReflect.Descriptor instead.
func (*GalacticKittensEnterGameResponse) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{8}
}

func (x *GalacticKittensEnterGameResponse) GetResult() *MessageResult {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *GalacticKittensEnterGameResponse) GetPlayer() *GalacticKittensPlayerInfo {
	if x != nil {
		return x.Player
	}
	return nil
}

// 游戏完成 内部
type GalacticKittensGameFinishRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId int64 `protobuf:"varint,1,opt,name=playerId,proto3" json:"playerId,omitempty"` //玩家id
}

func (x *GalacticKittensGameFinishRequest) Reset() {
	*x = GalacticKittensGameFinishRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensGameFinishRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensGameFinishRequest) ProtoMessage() {}

func (x *GalacticKittensGameFinishRequest) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensGameFinishRequest.ProtoReflect.Descriptor instead.
func (*GalacticKittensGameFinishRequest) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{9}
}

func (x *GalacticKittensGameFinishRequest) GetPlayerId() int64 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

// 游戏完成 内部
type GalacticKittensGameFinishResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *MessageResult           `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"` //错误信息
	Room   *GalacticKittensRoomInfo `protobuf:"bytes,2,opt,name=room,proto3" json:"room,omitempty"`     //房间信息
}

func (x *GalacticKittensGameFinishResponse) Reset() {
	*x = GalacticKittensGameFinishResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_game_galactic_kittens_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GalacticKittensGameFinishResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GalacticKittensGameFinishResponse) ProtoMessage() {}

func (x *GalacticKittensGameFinishResponse) ProtoReflect() protoreflect.Message {
	mi := &file_game_galactic_kittens_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GalacticKittensGameFinishResponse.ProtoReflect.Descriptor instead.
func (*GalacticKittensGameFinishResponse) Descriptor() ([]byte, []int) {
	return file_game_galactic_kittens_proto_rawDescGZIP(), []int{10}
}

func (x *GalacticKittensGameFinishResponse) GetResult() *MessageResult {
	if x != nil {
		return x.Result
	}
	return nil
}

func (x *GalacticKittensGameFinishResponse) GetRoom() *GalacticKittensRoomInfo {
	if x != nil {
		return x.Room
	}
	return nil
}

var File_game_galactic_kittens_proto protoreflect.FileDescriptor

var file_game_galactic_kittens_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x67, 0x61, 0x6d, 0x65, 0x5f, 0x67, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x5f,
	0x6b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x1f, 0x47,
	0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x65, 0x72, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4a, 0x0a, 0x20, 0x47, 0x61,
	0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74,
	0x65, 0x72, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26,
	0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x4f, 0x0a, 0x1f, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x04, 0x72, 0x6f, 0x6f,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x6d, 0x22, 0x73, 0x0a, 0x17, 0x47, 0x61, 0x6c, 0x61, 0x63,
	0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e,
	0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x32, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74,
	0x74, 0x65, 0x6e, 0x73, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0xc3, 0x01, 0x0a,
	0x19, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x69, 0x63, 0x6b, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x69, 0x63, 0x6b, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x65,
	0x70, 0x61, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x70, 0x6f,
	0x77, 0x65, 0x72, 0x55, 0x70, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0c, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x55, 0x70, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x68, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x68, 0x70, 0x12, 0x12,
	0x0a, 0x04, 0x69, 0x63, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x63,
	0x6f, 0x6e, 0x22, 0x39, 0x0a, 0x1d, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69,
	0x74, 0x74, 0x65, 0x6e, 0x73, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x70, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x22, 0x48, 0x0a,
	0x1e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x26, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x3d, 0x0a, 0x1f, 0x47, 0x61, 0x6c, 0x61, 0x63,
	0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x47,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x7e, 0x0a, 0x20, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x45, 0x6e, 0x74, 0x65, 0x72, 0x47, 0x61,
	0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74,
	0x74, 0x65, 0x6e, 0x73, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x22, 0x3e, 0x0a, 0x20, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x6c,
	0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x79, 0x0a, 0x21, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74,
	0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6e,
	0x69, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x72,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x06, 0x72, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x2c, 0x0a, 0x04, 0x72, 0x6f, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74,
	0x65, 0x6e, 0x73, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x72, 0x6f, 0x6f,
	0x6d, 0x32, 0x6e, 0x0a, 0x1a, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74,
	0x74, 0x65, 0x6e, 0x73, 0x47, 0x61, 0x6d, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x50, 0x0a, 0x09, 0x65, 0x6e, 0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x2e, 0x47,
	0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x45, 0x6e, 0x74, 0x65, 0x72, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x32, 0x72, 0x0a, 0x1b, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74,
	0x74, 0x65, 0x6e, 0x73, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x53, 0x0a, 0x0a, 0x67, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x12, 0x21,
	0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x73,
	0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x22, 0x2e, 0x47, 0x61, 0x6c, 0x61, 0x63, 0x74, 0x69, 0x63, 0x4b, 0x69, 0x74, 0x74,
	0x65, 0x6e, 0x73, 0x47, 0x61, 0x6d, 0x65, 0x46, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_game_galactic_kittens_proto_rawDescOnce sync.Once
	file_game_galactic_kittens_proto_rawDescData = file_game_galactic_kittens_proto_rawDesc
)

func file_game_galactic_kittens_proto_rawDescGZIP() []byte {
	file_game_galactic_kittens_proto_rawDescOnce.Do(func() {
		file_game_galactic_kittens_proto_rawDescData = protoimpl.X.CompressGZIP(file_game_galactic_kittens_proto_rawDescData)
	})
	return file_game_galactic_kittens_proto_rawDescData
}

var file_game_galactic_kittens_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_game_galactic_kittens_proto_goTypes = []interface{}{
	(*GalacticKittensEnterRoomRequest)(nil),   // 0: GalacticKittensEnterRoomRequest
	(*GalacticKittensEnterRoomResponse)(nil),  // 1: GalacticKittensEnterRoomResponse
	(*GalacticKittensRoomInfoResponse)(nil),   // 2: GalacticKittensRoomInfoResponse
	(*GalacticKittensRoomInfo)(nil),           // 3: GalacticKittensRoomInfo
	(*GalacticKittensPlayerInfo)(nil),         // 4: GalacticKittensPlayerInfo
	(*GalacticKittensPrepareRequest)(nil),     // 5: GalacticKittensPrepareRequest
	(*GalacticKittensPrepareResponse)(nil),    // 6: GalacticKittensPrepareResponse
	(*GalacticKittensEnterGameRequest)(nil),   // 7: GalacticKittensEnterGameRequest
	(*GalacticKittensEnterGameResponse)(nil),  // 8: GalacticKittensEnterGameResponse
	(*GalacticKittensGameFinishRequest)(nil),  // 9: GalacticKittensGameFinishRequest
	(*GalacticKittensGameFinishResponse)(nil), // 10: GalacticKittensGameFinishResponse
	(*MessageResult)(nil),                     // 11: MessageResult
}
var file_game_galactic_kittens_proto_depIdxs = []int32{
	11, // 0: GalacticKittensEnterRoomResponse.result:type_name -> MessageResult
	3,  // 1: GalacticKittensRoomInfoResponse.room:type_name -> GalacticKittensRoomInfo
	4,  // 2: GalacticKittensRoomInfo.player:type_name -> GalacticKittensPlayerInfo
	11, // 3: GalacticKittensPrepareResponse.result:type_name -> MessageResult
	11, // 4: GalacticKittensEnterGameResponse.result:type_name -> MessageResult
	4,  // 5: GalacticKittensEnterGameResponse.player:type_name -> GalacticKittensPlayerInfo
	11, // 6: GalacticKittensGameFinishResponse.result:type_name -> MessageResult
	3,  // 7: GalacticKittensGameFinishResponse.room:type_name -> GalacticKittensRoomInfo
	7,  // 8: GalacticKittensGameService.enterGame:input_type -> GalacticKittensEnterGameRequest
	9,  // 9: GalacticKittensMatchService.gameFinish:input_type -> GalacticKittensGameFinishRequest
	8,  // 10: GalacticKittensGameService.enterGame:output_type -> GalacticKittensEnterGameResponse
	10, // 11: GalacticKittensMatchService.gameFinish:output_type -> GalacticKittensGameFinishResponse
	10, // [10:12] is the sub-list for method output_type
	8,  // [8:10] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_game_galactic_kittens_proto_init() }
func file_game_galactic_kittens_proto_init() {
	if File_game_galactic_kittens_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_game_galactic_kittens_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensEnterRoomRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensEnterRoomResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensRoomInfoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensRoomInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensPlayerInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensPrepareRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensPrepareResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensEnterGameRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensEnterGameResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensGameFinishRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_game_galactic_kittens_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GalacticKittensGameFinishResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_game_galactic_kittens_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_game_galactic_kittens_proto_goTypes,
		DependencyIndexes: file_game_galactic_kittens_proto_depIdxs,
		MessageInfos:      file_game_galactic_kittens_proto_msgTypes,
	}.Build()
	File_game_galactic_kittens_proto = out.File
	file_game_galactic_kittens_proto_rawDesc = nil
	file_game_galactic_kittens_proto_goTypes = nil
	file_game_galactic_kittens_proto_depIdxs = nil
}
