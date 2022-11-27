package primaryAdapter

import (
	"context"

	"graft/pkg/infrastructure/adapter/clusterRpc"
	port "graft/pkg/infrastructure/port/primary"
)

type grpcApi struct {
	clusterRpc.UnimplementedClusterServer
	// Require redefining all rpc methods
	// because of forward compatibility
	adapter port.ServerAdapter
}

func NewGrpcApi(adapter port.ServerAdapter) *grpcApi {
	return &grpcApi{adapter: adapter}
}

func (s *grpcApi) AppendEntries(ctx context.Context, input *clusterRpc.AppendEntriesInput) (*clusterRpc.AppendEntriesOutput, error) {
	return s.adapter.AppendEntries(ctx, input)
}

func (s *grpcApi) RequestVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	return s.adapter.RequestVote(ctx, input)
}

func (s *grpcApi) PreVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	return s.adapter.PreVote(ctx, input)
}

func (s *grpcApi) Execute(ctx context.Context, input *clusterRpc.ExecuteInput) (*clusterRpc.ExecuteOutput, error) {
	return s.adapter.Execute(ctx, input)
}

func (s *grpcApi) LeadershipTransfer(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	return s.adapter.LeadershipTransfer(ctx, input)
}

func (s *grpcApi) Configuration(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.ConfigurationOutput, error) {
	return s.adapter.Configuration(ctx, input)
}

func (s *grpcApi) Shutdown(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	return s.adapter.Shutdown(ctx, input)
}

func (s *grpcApi) Ping(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	return s.adapter.Ping(ctx, input)
}
