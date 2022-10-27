package primaryAdapter

import (
	"context"
	"graft/app/infrastructure/adapter/rpc"
)

type repository interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}

type grpcApi struct {
	rpc.UnimplementedRpcServer
	repository
}

func NewGrpcApi(repository repository) *grpcApi {
	return &grpcApi{repository: repository}
}

func (s *grpcApi) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return s.repository.AppendEntries(ctx, input)
}

func (s *grpcApi) RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return s.repository.RequestVote(ctx, input)
}
