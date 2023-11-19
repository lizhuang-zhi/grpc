package main

import (
	"context"
	"fmt"
	"log"
	pb "rpc/client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PerRPCCredentials interface {
	GetRequestMetaData(ctx context.Context, uri ...string) (map[string]string, error) 
	RequireTransportSecurity() bool
}

type ClientTokenAuth struct {

}

func (c ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId": "leo",
		"appKey": "1123",
	}, nil
}

func (c ClientTokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	// 连接服务端, 此处禁用安全传输，没有加密和验证
	conn, err := grpc.Dial("127.0.0.1:9090", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 关闭连接

	// 建立连接
	client := pb.NewSayHelloClient(conn)
	// 调用服务端函数
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "Leo God"})
	fmt.Println(resp.GetResponseMsg())
}
