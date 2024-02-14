package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	svc "github.com/practicalgo/code/chap9/interceptors/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// 拦截器作用: 为客户端添加元数据
func metadataUnaryInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123-unary",
	)
	// 一元调用使用invoker调用服务端
	return invoker(ctxWithMetadata, method, req, reply, cc, opts...)
}

// 拦截器作用: 为客户端添加元数据
func metadataStreamInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	ctxWithMetadata := metadata.AppendToOutgoingContext(
		ctx,
		"Request-Id",
		"request-123-stream",
	)
	// 流调用使用streamer调用服务器
	clientStream, err := streamer(
		ctxWithMetadata,
		desc,
		cc,
		method,
		opts...,
	)
	return clientStream, err
}

func setupGrpcConn(addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		context.Background(),
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		// 客户端拦截器
		grpc.WithUnaryInterceptor(metadataUnaryInterceptor),   // 一元调用拦截器
		grpc.WithStreamInterceptor(metadataStreamInterceptor), // 流调用拦截器
	)
}

func getUserServiceClient(conn *grpc.ClientConn) svc.UsersClient {
	return svc.NewUsersClient(conn)
}

func getUser(
	client svc.UsersClient,
	u *svc.UserGetRequest,
) (*svc.UserGetReply, error) {
	return client.GetUser(context.Background(), u)
}

func setupChat(r io.Reader, w io.Writer, c svc.UsersClient) error {

	stream, err := c.GetHelp(context.Background())
	if err != nil {
		return err
	}
	for {
		scanner := bufio.NewScanner(r)
		prompt := "Request: "
		fmt.Fprint(w, prompt)

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return err
		}
		msg := scanner.Text()
		if msg == "quit" {
			break
		}
		request := svc.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&request)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("Response: %s\n", resp.Response)
	}
	return stream.CloseSend()
}

// 1. 进入该目录，执行 go build -o client 生成client可执行文件
// 2. 执行 ./client localhost:50051 GetUser
func main() {
	if len(os.Args) != 3 {
		log.Fatal(
			"Specify a gRPC server and method to call",
		)
	}
	serverAddr := os.Args[1] // 参数1
	methodName := os.Args[2] // 参数2

	conn, err := setupGrpcConn(serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := getUserServiceClient(conn)

	switch methodName {
	case "GetUser":
		result, err := getUser(
			c,
			&svc.UserGetRequest{Email: "jane@doe.com"},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(
			os.Stdout, "User: %s %s\n",
			result.User.FirstName,
			result.User.LastName,
		)
	case "GetHelp":
		err = setupChat(os.Stdin, os.Stdout, c)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unrecognized method name")
	}
}
