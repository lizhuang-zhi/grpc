package main

import (
	"context"
	"fmt"
	pb "grpc/protobuf/gen-pb"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10) // 建立一个内存通信信道
	s := grpc.NewServer()
	pb.RegisterGMServiceServer(s, &gmServer{})
	pb.RegisterRepoServer(s, &repoService{})

	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatal(err)
		}
	}()

	return s, l
}

func TestGMService(t *testing.T) {
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

// 测试服务端流
func TestRepoStreamService(t *testing.T) {
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

	repoClient := pb.NewRepoClient(client)
	stream, err := repoClient.GetRepos(context.Background(), &pb.RepoRequest{
		Id:        "repo-123",
		CreatorId: "user-123",
	})

	if err != nil {
		t.Fatal(err)
	}

	var repos []*pb.Repository
	for {
		repo, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		repos = append(repos, repo.Repo)
	}

	if len(repos) != 5 {
		t.Fatalf("expected 5 repos, got %d", len(repos))
	}

	for i, repo := range repos {
		gotRepoName := repo.Name
		expectedRepoName := fmt.Sprintf("repo-%d", i+1)

		if gotRepoName != expectedRepoName {
			t.Errorf("expected repo name %s, got %s", expectedRepoName, gotRepoName)
		}
	}

}
