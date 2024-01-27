package main

import (
	"context"
	"fmt"
	"log"
	pb "rpc/client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接服务端, 此处禁用安全传输，没有加密和验证
	conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭连接

	// 建立连接
	client := pb.NewSayHelloClient(conn)
	// 调用服务端函数
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Niu Bi"})
	ball, _ := client.PlayBall(context.Background(), &pb.Tools{Ball: "一个篮球", Count: 1})
	fmt.Println(resp.GetResponseMsg())
	fmt.Println(ball)
	fmt.Println(ball.GetMsg())
}
