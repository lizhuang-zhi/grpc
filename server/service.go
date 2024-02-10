package main

import (
	"context"
	"fmt"
	pb "grpc/protobuf/gen-pb"
	"grpc/server/common/gm"
	"grpc/server/mongo"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	statusConvert := status.Convert(err)
	if statusConvert.Code() != codes.OK {
		log.Fatalf("QueryPlayerInfo failed: %v - %v\n", statusConvert.Code(), statusConvert.Message())
		return nil, statusConvert.Err()
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

// repo server
type repoService struct {
	pb.UnimplementedRepoServer
}

func (s *repoService) GetRepos(in *pb.RepoRequest, stream pb.Repo_GetReposServer) error {
	log.Printf("Received request for repo with CreateId: %s Id: %s\n", in.CreatorId, in.Id)
	repo := pb.Repository{
		Id: in.Id,
		Owner: &pb.User{
			Id:        in.CreatorId,
			FirstName: "Jane",
		},
	}
	cnt := 1

	for {
		repo.Name = fmt.Sprintf("repo-%d", cnt)
		repo.Url = fmt.Sprintf(
			"https://git.emample.com/test/%s", repo.Name,
		)
		r := pb.RepoGetReply{
			Repo: &repo,
		}
		if err := stream.Send(&r); err != nil {
			return err
		}

		if cnt >= 5 {
			break
		}
		cnt++
	}
	return nil
}

func (s *repoService) CreateRepo(stream pb.Repo_CreateRepoServer) error {
	var repoContext *pb.RepoContext
	var data []byte

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Error(
				codes.Unknown,
				err.Error(),
			)
		}

		switch t := r.Body.(type) {
		case *pb.RepoCreateRequest_Context: // 读取上下文
			repoContext = r.GetContext()
		case *pb.RepoCreateRequest_Data: // 读取数据
			b := r.GetData()
			data = append(data, b...)
		case nil:
			return status.Error(
				codes.InvalidArgument,
				"Message doesn't contain context or data",
			)
		default:
			return status.Errorf(
				codes.FailedPrecondition,
				"Unexpected message type: %s",
				t,
			)
		}
	}

	repo := pb.Repository{
		Name: repoContext.Name,
		Url: fmt.Sprintf(
			"https://git.example.com/%s/%s",
			repoContext.CreatorId,
			repoContext.Name,
		),
	}

	r := pb.RepoCreateReply{
		Repo: &repo,
		Size: int32(len(data)),
	}

	return stream.SendAndClose(&r)
}

func main() {
	// 1. 开启端口
	listen, _ := net.Listen("tcp", ":9092")
	// 2. 创建grpc服务
	grpcServer := grpc.NewServer()
	// 3. 在grpc服务端注册服务
	pb.RegisterSayHelloServer(grpcServer, &server{})
	pb.RegisterGMServiceServer(grpcServer, &gmServer{})
	pb.RegisterRepoServer(grpcServer, &repoService{})
	// 4. 启动服务端
	grpcServer.Serve(listen)
}
