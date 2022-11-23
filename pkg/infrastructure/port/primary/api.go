package primaryPort

import (
	"context"

	"graft/pkg/domain"
	"graft/pkg/infrastructure/adapter/p2pRpc"
	"graft/pkg/services/rpc"
)

type ServerAdapter interface {
	AppendEntries(ctx context.Context, input *p2pRpc.AppendEntriesInput) (*p2pRpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
	PreVote(ctx context.Context, input *p2pRpc.RequestVoteInput) (*p2pRpc.RequestVoteOutput, error)
	//
	Execute(ctx context.Context, input *p2pRpc.ExecuteInput) (*p2pRpc.ExecuteOutput, error)
	ClusterConfiguration(ctx context.Context, input *p2pRpc.ClusterConfigurationInput) (*p2pRpc.ClusterConfigurationOutput, error)
}

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
		return nil, err
	}

	return &p2pRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}

func (p *rpcServerPort) Execute(ctx context.Context, input *p2pRpc.ExecuteInput) (*p2pRpc.ExecuteOutput, error) {
	_, err := p.adapter.Execute(&domain.ApiCommand{
		Type: domain.LogType(input.Type),
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (p *rpcServerPort) ClusterConfiguration(ctx context.Context, input *p2pRpc.ClusterConfigurationInput) (*p2pRpc.ClusterConfigurationOutput, error) {
	output, err := p.adapter.ClusterConfiguration()
	if err != nil {
		return nil, err
	}

	var peers []*p2pRpc.Peer
	for _, peer := range output.Peers {
		peers = append(peers, &p2pRpc.Peer{
			Id:     peer.Id,
			Host:   peer.Target(),
			Active: peer.Active,
		})
	}

	return &p2pRpc.ClusterConfigurationOutput{
		ElectionTimeout: uint32(output.ElectionTimeout),
		LeaderHeartbeat: uint32(output.LeaderHeartbeat),
		Peers:           peers,
	}, nil
}
