package main

import (
	"context"
	"fmt"
	pb "grpc/protobuf/gen-pb"
	"io"
	"log"
	"net"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startTestGrpcServer() (*grpc.Server, *bufconn.Listener) {
	l := bufconn.Listen(10) // 建立一个内存通信信道
	s := grpc.NewServer()
	pb.RegisterGMServiceServer(s, &gmServer{})
	pb.RegisterRepoServer(s, &repoService{})
	pb.RegisterUsersServer(s, &userService{})

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

// 测试客户端流
func TestRepoCreateStreamService(t *testing.T) {
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
	stream, err := repoClient.CreateRepo(
		context.Background(),
	)
	if err != nil {
		t.Fatal("CreateRepo", err)
	}

	c := pb.RepoCreateRequest_Context{
		Context: &pb.RepoContext{
			CreatorId: "user-123",
			Name:      "test-repo",
		},
	}

	r := pb.RepoCreateRequest{
		Body: &c,
	}

	err = stream.Send(&r)
	if err != nil {
		t.Fatal("StreamSend", err)
	}

	data := "Arbitrary Data Bytes"
	repoData := strings.NewReader(data)
	for {
		b, err := repoData.ReadByte()
		if err == io.EOF {
			break
		}

		bData := pb.RepoCreateRequest_Data{
			Data: []byte{b},
		}

		r := pb.RepoCreateRequest{
			Body: &bData,
		}

		err = stream.Send(&r)
		if err != nil {
			t.Fatal("StreamSend", err)
		}
		l.Close()
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		t.Fatal("CloseAndRecv", err)
	}

	expectedSize := int32(len(data))
	if resp.Size != expectedSize {
		t.Errorf(
			"Expected Repo Created to be: %d bytes Got back: %d",
			expectedSize,
			resp.Size,
		)
	}

	expectedRepoUrl := "https://git.example.com/user-123/test-repo"
	if resp.Repo.Url != expectedRepoUrl {
		t.Errorf(
			"Expected Repo URL to be: %s, Got: %s",
			expectedRepoUrl,
			resp.Repo.Url,
		)
	}
}

// 测试双向数据流
func TestGetHelp(t *testing.T) {
	_, l := startTestGrpcServer()

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

	userClient := pb.NewUsersClient(client)
	stream, err := userClient.GetHelp(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		msg := fmt.Sprintf("Hello, %d", i)
		r := pb.UserHelpRequest{
			Request: msg,
		}
		err := stream.Send(&r)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := stream.Recv()
		if err != nil {
			t.Fatal(err)
		}
		if resp.Response != msg {
			t.Errorf(
				"Expected Response to be: %s, Got: %s",
				msg,
				resp.Response,
			)
		}
	}
}
