// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: game_galactic_kittens.proto

package message

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GalacticKittensGameServiceClient is the client API for GalacticKittensGameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GalacticKittensGameServiceClient interface {
	// 进入游戏
	EnterGame(ctx context.Context, in *GalacticKittensEnterGameRequest, opts ...grpc.CallOption) (*GalacticKittensEnterGameResponse, error)
}

type galacticKittensGameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGalacticKittensGameServiceClient(cc grpc.ClientConnInterface) GalacticKittensGameServiceClient {
	return &galacticKittensGameServiceClient{cc}
}

func (c *galacticKittensGameServiceClient) EnterGame(ctx context.Context, in *GalacticKittensEnterGameRequest, opts ...grpc.CallOption) (*GalacticKittensEnterGameResponse, error) {
	out := new(GalacticKittensEnterGameResponse)
	err := c.cc.Invoke(ctx, "/GalacticKittensGameService/enterGame", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GalacticKittensGameServiceServer is the server API for GalacticKittensGameService service.
// All implementations must embed UnimplementedGalacticKittensGameServiceServer
// for forward compatibility
type GalacticKittensGameServiceServer interface {
	// 进入游戏
	EnterGame(context.Context, *GalacticKittensEnterGameRequest) (*GalacticKittensEnterGameResponse, error)
	mustEmbedUnimplementedGalacticKittensGameServiceServer()
}

// UnimplementedGalacticKittensGameServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGalacticKittensGameServiceServer struct {
}

func (UnimplementedGalacticKittensGameServiceServer) EnterGame(context.Context, *GalacticKittensEnterGameRequest) (*GalacticKittensEnterGameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnterGame not implemented")
}
func (UnimplementedGalacticKittensGameServiceServer) mustEmbedUnimplementedGalacticKittensGameServiceServer() {
}

// UnsafeGalacticKittensGameServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GalacticKittensGameServiceServer will
// result in compilation errors.
type UnsafeGalacticKittensGameServiceServer interface {
	mustEmbedUnimplementedGalacticKittensGameServiceServer()
}

func RegisterGalacticKittensGameServiceServer(s grpc.ServiceRegistrar, srv GalacticKittensGameServiceServer) {
	s.RegisterService(&GalacticKittensGameService_ServiceDesc, srv)
}

func _GalacticKittensGameService_EnterGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GalacticKittensEnterGameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GalacticKittensGameServiceServer).EnterGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GalacticKittensGameService/enterGame",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GalacticKittensGameServiceServer).EnterGame(ctx, req.(*GalacticKittensEnterGameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GalacticKittensGameService_ServiceDesc is the grpc.ServiceDesc for GalacticKittensGameService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GalacticKittensGameService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GalacticKittensGameService",
	HandlerType: (*GalacticKittensGameServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "enterGame",
			Handler:    _GalacticKittensGameService_EnterGame_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "game_galactic_kittens.proto",
}

// GalacticKittensMatchServiceClient is the client API for GalacticKittensMatchService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GalacticKittensMatchServiceClient interface {
	// 进入游戏
	GameFinish(ctx context.Context, in *GalacticKittensGameFinishRequest, opts ...grpc.CallOption) (*GalacticKittensGameFinishResponse, error)
}

type galacticKittensMatchServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGalacticKittensMatchServiceClient(cc grpc.ClientConnInterface) GalacticKittensMatchServiceClient {
	return &galacticKittensMatchServiceClient{cc}
}

func (c *galacticKittensMatchServiceClient) GameFinish(ctx context.Context, in *GalacticKittensGameFinishRequest, opts ...grpc.CallOption) (*GalacticKittensGameFinishResponse, error) {
	out := new(GalacticKittensGameFinishResponse)
	err := c.cc.Invoke(ctx, "/GalacticKittensMatchService/gameFinish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GalacticKittensMatchServiceServer is the server API for GalacticKittensMatchService service.
// All implementations must embed UnimplementedGalacticKittensMatchServiceServer
// for forward compatibility
type GalacticKittensMatchServiceServer interface {
	// 进入游戏
	GameFinish(context.Context, *GalacticKittensGameFinishRequest) (*GalacticKittensGameFinishResponse, error)
	mustEmbedUnimplementedGalacticKittensMatchServiceServer()
}

// UnimplementedGalacticKittensMatchServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGalacticKittensMatchServiceServer struct {
}

func (UnimplementedGalacticKittensMatchServiceServer) GameFinish(context.Context, *GalacticKittensGameFinishRequest) (*GalacticKittensGameFinishResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GameFinish not implemented")
}
func (UnimplementedGalacticKittensMatchServiceServer) mustEmbedUnimplementedGalacticKittensMatchServiceServer() {
}

// UnsafeGalacticKittensMatchServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GalacticKittensMatchServiceServer will
// result in compilation errors.
type UnsafeGalacticKittensMatchServiceServer interface {
	mustEmbedUnimplementedGalacticKittensMatchServiceServer()
}

func RegisterGalacticKittensMatchServiceServer(s grpc.ServiceRegistrar, srv GalacticKittensMatchServiceServer) {
	s.RegisterService(&GalacticKittensMatchService_ServiceDesc, srv)
}

func _GalacticKittensMatchService_GameFinish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GalacticKittensGameFinishRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GalacticKittensMatchServiceServer).GameFinish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/GalacticKittensMatchService/gameFinish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GalacticKittensMatchServiceServer).GameFinish(ctx, req.(*GalacticKittensGameFinishRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GalacticKittensMatchService_ServiceDesc is the grpc.ServiceDesc for GalacticKittensMatchService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GalacticKittensMatchService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GalacticKittensMatchService",
	HandlerType: (*GalacticKittensMatchServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "gameFinish",
			Handler:    _GalacticKittensMatchService_GameFinish_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "game_galactic_kittens.proto",
}