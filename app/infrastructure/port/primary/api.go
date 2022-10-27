package primaryPort

import (
	"context"
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/rpc"
	"graft/app/usecase/receiver"
)

type rpcServerPort struct {
	adapter receiver.UseCase
}

func NewRpcServerPort(adapter receiver.UseCase) *rpcServerPort {
	return &rpcServerPort{adapter}
}

func (p *rpcServerPort) AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error) {
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

func (p *rpcServerPort) RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error) {
	output, err := p.adapter.RequestVote(&entity.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})

	if err != nil {
		return nil, err
	}

	return &rpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}
