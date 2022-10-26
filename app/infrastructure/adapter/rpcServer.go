package adapter

import (
	"context"
	"fmt"
	"graft/app/infrastructure/adapter/rpc"
	"log"
	"net"

	"google.golang.org/grpc"
)

type repository interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequesVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}

type rpcApi struct {
	rpc.UnimplementedRpcServer
	repo repository
}

func (s *rpcApi) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return s.repo.AppendEntries(ctx, input)
}

func (s *rpcApi) RequesVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return s.repo.RequesVote(ctx, input)
}

type rpcServer struct {
	rpcApi
}

func NewRpcServer(repo repository) *rpcServer {
	return &rpcServer{rpcApi{repo: repo}}
}

func (r *rpcServer) Start(port string) {
	log.Println("START RECEIVER SERVER")

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: \n\t%v\n", err)
	}

	server := grpc.NewServer()
	rpc.RegisterRpcServer(server, &r.rpcApi)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: \n\t%v\n", err)
	}
}
