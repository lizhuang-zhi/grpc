// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: gm.proto

package service

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

const (
	GMService_ExcuteGM_FullMethodName = "/GMService/ExcuteGM"
)

// GMServiceClient is the client API for GMService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GMServiceClient interface {
	// rpc 服务函数名(参数) 返回 (返回参数)
	ExcuteGM(ctx context.Context, in *GMRequest, opts ...grpc.CallOption) (*GMResponse, error)
}

type gMServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGMServiceClient(cc grpc.ClientConnInterface) GMServiceClient {
	return &gMServiceClient{cc}
}

func (c *gMServiceClient) ExcuteGM(ctx context.Context, in *GMRequest, opts ...grpc.CallOption) (*GMResponse, error) {
	out := new(GMResponse)
	err := c.cc.Invoke(ctx, GMService_ExcuteGM_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GMServiceServer is the server API for GMService service.
// All implementations must embed UnimplementedGMServiceServer
// for forward compatibility
type GMServiceServer interface {
	// rpc 服务函数名(参数) 返回 (返回参数)
	ExcuteGM(context.Context, *GMRequest) (*GMResponse, error)
	mustEmbedUnimplementedGMServiceServer()
}

// UnimplementedGMServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGMServiceServer struct {
}

func (UnimplementedGMServiceServer) ExcuteGM(context.Context, *GMRequest) (*GMResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExcuteGM not implemented")
}
func (UnimplementedGMServiceServer) mustEmbedUnimplementedGMServiceServer() {}

// UnsafeGMServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GMServiceServer will
// result in compilation errors.
type UnsafeGMServiceServer interface {
	mustEmbedUnimplementedGMServiceServer()
}

func RegisterGMServiceServer(s grpc.ServiceRegistrar, srv GMServiceServer) {
	s.RegisterService(&GMService_ServiceDesc, srv)
}

func _GMService_ExcuteGM_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GMRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GMServiceServer).ExcuteGM(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GMService_ExcuteGM_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GMServiceServer).ExcuteGM(ctx, req.(*GMRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GMService_ServiceDesc is the grpc.ServiceDesc for GMService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GMService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "GMService",
	HandlerType: (*GMServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExcuteGM",
			Handler:    _GMService_ExcuteGM_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gm.proto",
}