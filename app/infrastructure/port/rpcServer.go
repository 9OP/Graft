package port

import (
	"context"
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/rpc"
)

// repository in adapter rpcServer
//
// type repository interface {
// 	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
// 	RequesVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
// }

type rpcServerAdapter interface {
	AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}

type rpcServerPort struct {
	adapter rpcServerAdapter
}

func NewRpcServerPort(adapter rpcServerAdapter) *rpcServerPort {
	return &rpcServerPort{adapter}
}

func (p *rpcServerPort) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
	/*
		rpc server has its raft rpc handlers listenning for triggering logic on server:
		1- receive incomming RPC from outside (*rpc...)
		2- convert input into relevant entity.rpc
		3- run entity through use case
		4- reconvert use case response into *rpc...
	*/
	output, err := p.adapter.AppendEntries(&entity.AppendEntriesInput{
		Term:         input.Term,
		LeaderId:     input.LeaderId,
		PrevLogIndex: input.PrevLogIndex,
		PrevLogTerm:  input.PrevLogTerm,
		Entries:      input.Entries,
		LeaderCommit: input.LeaderCommit,
	})

	if err != nil {
		return nil, err
	}

	return &rpc.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcServerPort) RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	return nil, nil
}
