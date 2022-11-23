package primaryAdapter

import (
	"context"

	"graft/pkg/infrastructure/adapter/p2pRpc"
)

// Clean =, try to remove repository, I am not sure this is really required

type repository interface {
	AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
	PreVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
	//
	Execute(ctx context.Context, input *p2pRpc.ExecuteInput) (*p2pRpc.ExecuteOutput, error)
	ClusterConfiguration(ctx context.Context, input *p2pRpc.ClusterConfigurationInput) (*p2pRpc.ClusterConfigurationOutput, error)
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

func (s *grpcApi) PreVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	return s.repository.PreVote(ctx, input)
}

func (s *grpcApi) Execute(ctx context.Context, input *p2pRpc.ExecuteInput) (*p2pRpc.ExecuteOutput, error) {
	return s.repository.Execute(ctx, input)
}

func (s *grpcApi) ClusterConfiguration(ctx context.Context, input *p2pRpc.ClusterConfigurationInput) (*p2pRpc.ClusterConfigurationOutput, error) {
	return s.repository.ClusterConfiguration(ctx, input)
}
