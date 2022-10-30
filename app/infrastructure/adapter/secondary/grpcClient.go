package secondaryAdapter

import (
	"context"
	"graft/app/infrastructure/adapter/p2pRpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UseCaseGrpcClient interface {
	AppendEntries(target string, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error)
	RequestVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
}

type grpcClient struct{}

func NewGrpcClient() *grpcClient {
	return &grpcClient{}
}

func withClient[K *p2pRpc.AppendEntriesOutput | *p2pRpc.RequestVoteOutput](target string, fn func(c p2pRpc.RpcClient) (K, error)) (K, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(target, creds)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := p2pRpc.NewRpcClient(conn)
	return fn(c)
}

func (r *grpcClient) AppendEntries(target string, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	return withClient(
		target,
		func(c p2pRpc.RpcClient) (*p2pRpc.AppendEntriesOutput, error) {
			return c.AppendEntries(context.Background(), input)
		})
}

func (r *grpcClient) RequestVote(target string, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	return withClient(
		target,
		func(c p2pRpc.RpcClient) (*p2pRpc.RequestVoteOutput, error) {
			return c.RequestVote(context.Background(), input)
		})
}
