package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	pb "rpc/server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// hello server
type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("Hello " + req.RequestName)

	// 获取元数据的信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("未传入token")
	}

	var appId string
	var appKey string
	if v, ok := md["appId"]; ok {
		appId = v[0]
	}
	if v, ok := md["appKey"]; ok {
		appKey = v[0]
	}

	if appId != "leo" || appKey != "123" {
		fmt.Println("token不正确")
		return nil, errors.New("token不正确")
	}

	return &pb.HelloResponse{ResponseMsg: "Hello " + req.RequestName}, nil
}

func main() {
	// 1. 开启端口
	listen, _ := net.Listen("tcp", ":9090")
	// 2. 创建grpc服务
	grpcServer := grpc.NewServer()
	// 3. 在grpc服务端注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	// 4. 启动服务端
	grpcServer.Serve(listen)
}
