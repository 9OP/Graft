package adapter

import (
	"context"
	"graft/app/infrastructure/adapter/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// This adapter works the same as a plain database driver

type rpcClient struct{}

func NewRpcClient() *rpcClient {
	return &rpcClient{}
}

func withClient[K *rpc.AppendEntriesOutput | *rpc.RequestVoteOutput](target string, fn func(c rpc.RpcClient) (K, error)) (K, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(target, creds)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpc.NewRpcClient(conn)
	return fn(c)
}

func (r *rpcClient) AppendEntries(target string, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return withClient(
		target,
		func(c rpc.RpcClient) (*rpc.AppendEntriesOutput, error) {
			return c.AppendEntries(context.Background(), input)
		})
}

func (r *rpcClient) RequestVote(target string, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return withClient(
		target,
		func(c rpc.RpcClient) (*rpc.RequestVoteOutput, error) {
			return c.RequestVote(context.Background(), input)
		})
}
