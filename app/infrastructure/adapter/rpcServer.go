package adapter

import (
	"context"
	"graft/app/infrastructure/adapter/rpc"
)

type repository interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}

type RpcApi struct {
	rpc.UnimplementedRpcServer
	repo repository
}

func NewRpcApi(repo repository) *RpcApi {
	return &RpcApi{repo: repo}
}

func (s *RpcApi) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return s.repo.AppendEntries(ctx, input)
}

func (s *RpcApi) RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return s.repo.RequestVote(ctx, input)
}
