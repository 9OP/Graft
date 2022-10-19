package clients

import (
	"context"
	"fmt"
	"graft/src2/entity"
	"graft/src2/rpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Runner struct{}

func withClient[K *rpc.AppendEntriesOutput | *rpc.RequestVoteOutput](peer entity.Peer, fn func(c rpc.RpcClient) (K, error)) (K, error) {
	target := fmt.Sprintf("%s:%s", peer.Host, peer.Port)
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(target, creds)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	c := rpc.NewRpcClient(conn)
	return fn(c)
}

func (c *Runner) AppendEntries(peer entity.Peer, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	return withClient(
		peer,
		func(c rpc.RpcClient) (*rpc.AppendEntriesOutput, error) {
			return c.AppendEntries(context.Background(), input)
		})
}

func (c *Runner) RequestVote(peer entity.Peer, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	return withClient(
		peer,
		func(c rpc.RpcClient) (*rpc.RequestVoteOutput, error) {
			return c.RequestVote(context.Background(), input)
		})
}
