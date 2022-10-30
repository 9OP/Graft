package secondaryPort

import (
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/p2pRpc"
	adapter "graft/app/infrastructure/adapter/secondary"
)

type rpcClientPort struct {
	adapter adapter.UseCaseGrpcClient
}

func NewRpcClientPort(adapter adapter.UseCaseGrpcClient) *rpcClientPort {
	return &rpcClientPort{adapter}
}

func (p *rpcClientPort) AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	target := peer.Target()
	entries := make([]*p2pRpc.LogEntry, 0, len(input.Entries))
	for _, log := range input.Entries {
		entries = append(entries, &p2pRpc.LogEntry{Term: log.Term, Value: log.Value})
	}
	output, err := p.adapter.AppendEntries(target, &p2pRpc.AppendEntriesInput{
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

	return &entity.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcClientPort) RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	target := peer.Target()
	output, err := p.adapter.RequestVote(target, &p2pRpc.RequestVoteInput{
		Term:         input.Term,
		CandidateId:  input.CandidateId,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})

	if err != nil {
		return nil, err
	}

	return &entity.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}
