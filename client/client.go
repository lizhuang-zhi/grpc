package main

import (
	"context"
	"fmt"
	"log"

	pb "grpc/protobuf/gen-pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接服务端, 此处禁用安全传输，没有加密和验证
	/*
		grpc.WithBlock(): 阻塞连接，直到连接成功或者超时
	*/
	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭连接

	// 建立连接
	client := getSayHelloServiceClient(conn)
	gmClient := getGMServiceClient(conn)
	// 调用服务端函数
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Niu Bi"})
	ball, _ := client.PlayBall(context.Background(), &pb.Tools{Ball: "一个篮球", Count: 1})
	fmt.Println(resp.GetResponseMsg())
	fmt.Println(ball)
	fmt.Println(ball.GetMsg())

	// 调用gm服务端函数
	gmResp, _ := gmClient.ExcuteGM(context.Background(), &pb.GMRequest{
		Command: "additem",
		Args:    "1001,1002,1003",
		// PlayerID: "leo666",
		PlayerID: "", // 报错
	})
	fmt.Println(gmResp)
	fmt.Println(gmResp.GetMsg())
	fmt.Println(gmResp.GetCode())
}

func getSayHelloServiceClient(conn *grpc.ClientConn) pb.SayHelloClient {
	return pb.NewSayHelloClient(conn)
}

func getGMServiceClient(conn *grpc.ClientConn) pb.GMServiceClient {
	return pb.NewGMServiceClient(conn)
}
