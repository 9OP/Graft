package primaryAdapter

import (
	"context"
	"graft/pkg/infrastructure/adapter/p2pRpc"
)

type repository interface {
	AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
}

type grpcApi struct {
	p2pRpc.UnimplementedRpcServer
	repository
}

func NewGrpcApi(repository repository) *grpcApi {
	return &grpcApi{repository: repository}
}

func (s *grpcApi) AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	return s.repository.AppendEntries(ctx, input)
}

func (s *grpcApi) RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	return s.repository.RequestVote(ctx, input)
}
