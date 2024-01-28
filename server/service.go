package main

import (
	"context"
	"net"
	"rpc/server/common/gm"
	"rpc/server/mongo"
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
		Site:   "篮球场地",
		Msg:    tools.Ball,
	}, nil
}

// gm server
type gmServer struct {
	pb.UnimplementedGMServiceServer
}

func (s *gmServer) ExcuteGM(ctx context.Context, req *pb.GMRequest) (*pb.GMResponse, error) {
	// 获取命令
	commnd := req.Command
	// 获取参数
	args := req.Args
	// 获取玩家id
	playerID := req.PlayerID

	// 查询玩家信息
	playerInfo, err := mongo.QueryPlayerInfo(playerID)
	if err != nil {
		return nil, err
	}

	if playerInfo.Level < 2 {
		return nil, nil
	}

	// 执行命令
	res, err := gm.ExecuteGMCommand(commnd, args, playerID)
	if err != nil {
		return nil, err
	}

	return &pb.GMResponse{
		Code: 0,
		Msg:  res,
	}, nil
}

func main() {
	// 1. 开启端口
	listen, _ := net.Listen("tcp", ":9092")
	// 2. 创建grpc服务
	grpcServer := grpc.NewServer()
	// 3. 在grpc服务端注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	pb.RegisterGMServiceServer(grpcServer, &gmServer{})
	// 4. 启动服务端
	grpcServer.Serve(listen)
}
