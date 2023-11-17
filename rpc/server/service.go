package main

import (
	"context"
	"fmt"
	"net"
	pb "rpc/server/proto"

	"google.golang.org/grpc"
)

// hello server
type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("Hello " + req.RequestName)
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
