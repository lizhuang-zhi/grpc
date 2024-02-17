package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	svc "github.com/practicalgo/code/chap9/interceptor-chain/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userService struct {
	svc.UnimplementedUsersServer
}

func (s *userService) GetUser(
	ctx context.Context,
	in *svc.UserGetRequest,
) (*svc.UserGetReply, error) {

	log.Printf(
		"Received request for user with Email: %s Id: %s\n",
		in.Email,
		in.Id,
	)
	components := strings.Split(in.Email, "@")
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := svc.User{
		Id:        in.Id,
		FirstName: components[0],
		LastName:  components[1],
		Age:       36,
	}
	return &svc.UserGetReply{User: &u}, nil
}

func (s *userService) GetHelp(
	stream svc.Users_GetHelpServer,
) error {
	for {

		request, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Request receieved: %s\n", request.Request)
		response := svc.UserHelpReply{
			Response: request.Request,
		}
		err = stream.Send(&response)
		if err != nil {
			return err
		}
	}
	return nil
}

func registerServices(s *grpc.Server) {
	svc.RegisterUsersServer(s, &userService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}

func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(
		// 设置拦截器链
		grpc.ChainUnaryInterceptor(
			metricUnaryInterceptor, // 链接多个拦截器
			loggingUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			metricStreamInterceptor,
			loggingStreamInterceptor,
		),
	)
	registerServices(s)
	log.Fatal(startServer(s, lis))
}

// 服务端包装流
type wrappedServerStream struct {
	grpc.ServerStream
}

func (s wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send msg called: %T", m)
	return s.ServerStream.SendMsg(m)
}

func (s wrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("Waiting to receive a message: %T", m)
	return s.ServerStream.RecvMsg(m)
}

func loggingUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return resp, err
}

func loggingStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	// ***** 服务端包装流 *****
	// 将原来的(interceptor目录下)stream包装成wrappedServerStream
	serverStream := wrappedServerStream{
		ServerStream: stream,
	}
	err := handler(srv, serverStream)
	// ***** 服务端包装流 *****

	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Print("No metadata")
	}
	log.Printf("Method:%s, Error:%v, Request-Id:%s",
		info.FullMethod,
		err,
		md.Get("Request-Id"),
	)
	return err
}

// 性能监控的拦截器
func metricUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return resp, err
}

func metricStreamInterceptor(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {

	start := time.Now()
	err := handler(srv, stream)
	end := time.Now()
	log.Printf("Method:%s, Duration:%s",
		info.FullMethod,
		end.Sub(start),
	)
	return err
}
