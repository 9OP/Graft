package primaryPort

import (
	"context"
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/p2pRpc"
	"graft/app/usecase/receiver"
)

type rpcServerPort struct {
	adapter receiver.UseCase
}

func NewRpcServerPort(adapter receiver.UseCase) *rpcServerPort {
	return &rpcServerPort{adapter}
}

func (p *rpcServerPort) AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	entries := make([]entity.LogEntry, 0, len(input.Entries))
	for _, log := range input.Entries {
		entries = append(entries, entity.LogEntry{Term: log.Term, Value: log.Value})
	}
	output, err := p.adapter.AppendEntries(&entity.AppendEntriesInput{
		Term:         input.Term,
		LeaderId:     input.LeaderId,
		PrevLogIndex: input.PrevLogIndex,
		PrevLogTerm:  input.PrevLogTerm,
		Entries:      entries,
		LeaderCommit: input.LeaderCommit,
	})

	if err != nil {
		return nil, err
	}

	return &p2pRpc.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcServerPort) RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	output, err := p.adapter.RequestVote(&entity.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})

	if err != nil {
		return nil, err
	}

	return &p2pRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}
