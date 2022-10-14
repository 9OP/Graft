package main

import (
	"context"
	"fmt"
	"graft/src/rpc"

	"google.golang.org/grpc"
)

type RpcClient struct {
	host string
	port string
}

type rpcFunc = func(c rpc.RpcClient) (interface{}, error)

func (clt *RpcClient) withClient(fn rpcFunc) (interface{}, error) {
	target := fmt.Sprintf("%s:%s", clt.host, clt.port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()
	c := rpc.NewRpcClient(conn)
	return fn(c)
}

func (clt *RpcClient) SendRequestVoteRpc(input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	fn := func(c rpc.RpcClient) (interface{}, error) {
		return c.RequestVote(context.Background(), input)
	}
	res, err := clt.withClient(fn)
	if err != nil {
		return nil, err
	}
	return res.(*rpc.RequestVoteOutput), nil
}

func (clt *RpcClient) SendAppendEntriesRpc(input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	fn := func(c rpc.RpcClient) (interface{}, error) {
		return c.AppendEntries(context.Background(), input)
	}
	res, err := clt.withClient(fn)
	if err != nil {
		return nil, err
	}
	return res.(*rpc.AppendEntriesOutput), err
}
