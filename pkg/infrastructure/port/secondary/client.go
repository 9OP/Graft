package secondaryPort

import (
	"graft/pkg/domain"
	"graft/pkg/infrastructure/adapter/p2pRpc"
	adapter "graft/pkg/infrastructure/adapter/secondary"
)

type rpcClientPort struct {
	adapter adapter.UseCaseGrpcClient
}

func NewRpcClientPort(adapter adapter.UseCaseGrpcClient) *rpcClientPort {
	return &rpcClientPort{adapter}
}

func (p *rpcClientPort) AppendEntries(peer domain.Peer, input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
	target := peer.TargetP2p()
	entries := make([]*p2pRpc.LogEntry, 0, len(input.Entries))
	for _, log := range input.Entries {
		logType := p2pRpc.LogEntry_LogType(p2pRpc.LogEntry_LogType_value[log.Type.String()])
		entry := &p2pRpc.LogEntry{
			Index: log.Index,
			Term:  log.Term,
			Data:  log.Data,
			Type:  logType,
		}
		entries = append(entries, entry)
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

	return &domain.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcClientPort) RequestVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	target := peer.TargetP2p()
	output, err := p.adapter.RequestVote(target, &p2pRpc.RequestVoteInput{
		Term:         input.Term,
		CandidateId:  input.CandidateId,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		return nil, err
	}

	return &domain.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}

func (p *rpcClientPort) PreVote(peer domain.Peer, input *domain.RequestVoteInput) (*domain.RequestVoteOutput, error) {
	target := peer.TargetP2p()
	output, err := p.adapter.PreVote(target, &p2pRpc.RequestVoteInput{
		Term:         input.Term,
		CandidateId:  input.CandidateId,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	return &domain.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}
