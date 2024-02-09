package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "grpc/protobuf/gen-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type dummyGMServiceServer struct {
	pb.UnimplementedGMServiceServer
}

func (s *dummyGMServiceServer) ExcuteGM(
	ctx context.Context,
	req *pb.GMRequest,
) (*pb.GMResponse, error) {
	return &pb.GMResponse{
		Code: 0,
		Msg:  "success",
	}, nil
}

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10) // 建立一个内存通信信道
	s := grpc.NewServer()
	pb.RegisterGMServiceServer(s, &dummyGMServiceServer{})

	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	return s, l
}

func TestExcuteGM(t *testing.T) {
	s, l := startTestGrpcServer()
	defer s.GracefulStop() // 优雅关闭服务

	// 创建一个拨号器
	bufconnDialer := func(ctx context.Context, address string) (net.Conn, error) {
		return l.Dial()
	}

	// 创建特殊配置客户端
	client, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(bufconnDialer),
	)
	if err != nil {
		t.Fatal(err)
	}

	GMClient := pb.NewGMServiceClient(client)

	resp, err := GMClient.ExcuteGM(
		context.Background(),
		&pb.GMRequest{
			Command:  "additem",
			Args:     "1001,1002,1003",
			PlayerID: "leo666",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	if resp.GetCode() != 0 {
		t.Fatal("error", resp.GetMsg())
	}
}
