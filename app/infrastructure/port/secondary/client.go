package secondaryPort

import (
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/rpc"
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
	output, err := p.adapter.AppendEntries(target, &rpc.AppendEntriesInput{
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

	return &entity.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcClientPort) RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error) {
	target := peer.Target()
	output, err := p.adapter.RequestVote(target, &rpc.RequestVoteInput{
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
