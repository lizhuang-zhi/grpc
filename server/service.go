package main

import (
	"context"
	"net"
	pb "rpc/server/proto"

	"google.golang.org/grpc"
)

// hello server
type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{ResponseMsg: "Hello " + req.RequestName}, nil
}

	// PlayBall(context.Context, *Tools) (*PlayBallStatus, error)
func (s *server) PlayBall(ctx context.Context, tools *pb.Tools) (*pb.PlayBallStatus, error) {
	return &pb.PlayBallStatus{
		People: "Leo",
		Site: "篮球场地",
		Msg: tools.Ball,
	}, nil
}

func main() {
	// 1. 开启端口
	listen, _ := net.Listen("tcp", ":9092")
	// 2. 创建grpc服务
	grpcServer := grpc.NewServer()
	// 3. 在grpc服务端注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	// 4. 启动服务端
	grpcServer.Serve(listen)
}
