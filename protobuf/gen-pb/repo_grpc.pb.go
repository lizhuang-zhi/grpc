// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.2
// source: protobuf/proto/repo.proto

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
	Repo_GetRepos_FullMethodName   = "/Repo/GetRepos"
	Repo_CreateRepo_FullMethodName = "/Repo/CreateRepo"
)

// RepoClient is the client API for Repo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RepoClient interface {
	// 返回stream流
	GetRepos(ctx context.Context, in *RepoRequest, opts ...grpc.CallOption) (Repo_GetReposClient, error)
	CreateRepo(ctx context.Context, opts ...grpc.CallOption) (Repo_CreateRepoClient, error)
}

type repoClient struct {
	cc grpc.ClientConnInterface
}

func NewRepoClient(cc grpc.ClientConnInterface) RepoClient {
	return &repoClient{cc}
}

func (c *repoClient) GetRepos(ctx context.Context, in *RepoRequest, opts ...grpc.CallOption) (Repo_GetReposClient, error) {
	stream, err := c.cc.NewStream(ctx, &Repo_ServiceDesc.Streams[0], Repo_GetRepos_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &repoGetReposClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Repo_GetReposClient interface {
	Recv() (*RepoGetReply, error)
	grpc.ClientStream
}

type repoGetReposClient struct {
	grpc.ClientStream
}

func (x *repoGetReposClient) Recv() (*RepoGetReply, error) {
	m := new(RepoGetReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *repoClient) CreateRepo(ctx context.Context, opts ...grpc.CallOption) (Repo_CreateRepoClient, error) {
	stream, err := c.cc.NewStream(ctx, &Repo_ServiceDesc.Streams[1], Repo_CreateRepo_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &repoCreateRepoClient{stream}
	return x, nil
}

type Repo_CreateRepoClient interface {
	Send(*RepoCreateRequest) error
	CloseAndRecv() (*RepoCreateReply, error)
	grpc.ClientStream
}

type repoCreateRepoClient struct {
	grpc.ClientStream
}

func (x *repoCreateRepoClient) Send(m *RepoCreateRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *repoCreateRepoClient) CloseAndRecv() (*RepoCreateReply, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(RepoCreateReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RepoServer is the server API for Repo service.
// All implementations must embed UnimplementedRepoServer
// for forward compatibility
type RepoServer interface {
	// 返回stream流
	GetRepos(*RepoRequest, Repo_GetReposServer) error
	CreateRepo(Repo_CreateRepoServer) error
	mustEmbedUnimplementedRepoServer()
}

// UnimplementedRepoServer must be embedded to have forward compatible implementations.
type UnimplementedRepoServer struct {
}

func (UnimplementedRepoServer) GetRepos(*RepoRequest, Repo_GetReposServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRepos not implemented")
}
func (UnimplementedRepoServer) CreateRepo(Repo_CreateRepoServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateRepo not implemented")
}
func (UnimplementedRepoServer) mustEmbedUnimplementedRepoServer() {}

// UnsafeRepoServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RepoServer will
// result in compilation errors.
type UnsafeRepoServer interface {
	mustEmbedUnimplementedRepoServer()
}

func RegisterRepoServer(s grpc.ServiceRegistrar, srv RepoServer) {
	s.RegisterService(&Repo_ServiceDesc, srv)
}

func _Repo_GetRepos_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RepoRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RepoServer).GetRepos(m, &repoGetReposServer{stream})
}

type Repo_GetReposServer interface {
	Send(*RepoGetReply) error
	grpc.ServerStream
}

type repoGetReposServer struct {
	grpc.ServerStream
}

func (x *repoGetReposServer) Send(m *RepoGetReply) error {
	return x.ServerStream.SendMsg(m)
}

func _Repo_CreateRepo_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RepoServer).CreateRepo(&repoCreateRepoServer{stream})
}

type Repo_CreateRepoServer interface {
	SendAndClose(*RepoCreateReply) error
	Recv() (*RepoCreateRequest, error)
	grpc.ServerStream
}

type repoCreateRepoServer struct {
	grpc.ServerStream
}

func (x *repoCreateRepoServer) SendAndClose(m *RepoCreateReply) error {
	return x.ServerStream.SendMsg(m)
}

func (x *repoCreateRepoServer) Recv() (*RepoCreateRequest, error) {
	m := new(RepoCreateRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Repo_ServiceDesc is the grpc.ServiceDesc for Repo service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Repo_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Repo",
	HandlerType: (*RepoServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRepos",
			Handler:       _Repo_GetRepos_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CreateRepo",
			Handler:       _Repo_CreateRepo_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/proto/repo.proto",
}
