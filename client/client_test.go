package main

import (
	"context"
	"log"
	"net"
	"testing"

	pb "grpc/protobuf/gen-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/encoding/protojson"
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

// 测试客户端
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

// protojson: 将json字符串映射到消息类型对象
func TestTransJson(t *testing.T) {
	// 模拟json字符串数据
	jsonString := `{"command": "additem", "args": "1001,1002,1003", "playerID": "leo666"}`
	gmReq, err := createGMRequest(jsonString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gmReq)
}

// json字符串转结构体
func createGMRequest(jsonQuery string) (*pb.GMRequest, error) {
	gm_request := pb.GMRequest{}
	err := protojson.Unmarshal([]byte(jsonQuery), &gm_request)
	if err != nil {
		return nil, err
	}
	return &gm_request, nil
}
