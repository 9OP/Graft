package primaryAdapter

import (
	"context"

	"graft/pkg/infrastructure/adapter/p2pRpc"
	port "graft/pkg/infrastructure/port/primary"
)

type grpcApi struct {
	p2pRpc.UnimplementedRpcServer
	// Require redefining all rpc methods
	// because of forward compatibility
	adapter port.ServerAdapter
}

func NewGrpcApi(adapter port.ServerAdapter) *grpcApi {
	return &grpcApi{adapter: adapter}
}

func (s *grpcApi) AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	return s.adapter.AppendEntries(ctx, input)
}

func (s *grpcApi) RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	return s.adapter.RequestVote(ctx, input)
}

func (s *grpcApi) PreVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	return s.adapter.PreVote(ctx, input)
}

func (s *grpcApi) Execute(ctx context.Context, input *p2pRpc.ExecuteInput) (*p2pRpc.ExecuteOutput, error) {
	return s.adapter.Execute(ctx, input)
}

func (s *grpcApi) ClusterConfiguration(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.ClusterConfigurationOutput, error) {
	return s.adapter.ClusterConfiguration(ctx, input)
}

func (s *grpcApi) Shutdown(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.Nil, error) {
	return s.adapter.Shutdown(ctx, input)
}
