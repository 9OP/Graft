package primaryPort

import (
	"context"

	"graft/pkg/domain"
	"graft/pkg/infrastructure/adapter/p2pRpc"
	"graft/pkg/services/rpc"

	log "github.com/sirupsen/logrus"
)

type rpcServerPort struct {
	adapter rpc.UseCase
}

func NewRpcServerPort(adapter rpc.UseCase) *rpcServerPort {
	return &rpcServerPort{adapter}
}

func (p *rpcServerPort) AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error) {
	entries := make([]domain.LogEntry, 0, len(input.Entries))
	for _, log := range input.Entries {
		logType := domain.LogType(log.Type)
		entry := domain.LogEntry{
			Term: log.Term,
			Data: log.Data,
			Type: logType,
		}
		entries = append(entries, entry)
	}
	output, err := p.adapter.AppendEntries(&domain.AppendEntriesInput{
		Term:         input.Term,
		LeaderId:     input.LeaderId,
		PrevLogIndex: input.PrevLogIndex,
		PrevLogTerm:  input.PrevLogTerm,
		Entries:      entries,
		LeaderCommit: input.LeaderCommit,
	})
	if err != nil {
		log.Errorf("RPC RESP APPEND_ENTRIES %s\n", input.LeaderId)
		return nil, err
	}

	return &p2pRpc.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcServerPort) RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	output, err := p.adapter.RequestVote(&domain.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		log.Errorf("RPC RESP REQUEST_VOTE %s\n", input.CandidateId)
		return nil, err
	}

	return &p2pRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}

func (p *rpcServerPort) PreVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error) {
	output, err := p.adapter.PreVote(&domain.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		log.Errorf("RPC RESP PRE_VOTE %s\n", input.CandidateId)
		return nil, err
	}

	return &p2pRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}
