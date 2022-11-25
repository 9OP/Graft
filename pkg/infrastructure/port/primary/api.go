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
	ClusterConfiguration(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.ClusterConfigurationOutput, error)
	Shutdown(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.Nil, error)
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
			Index: log.Index,
			Term:  log.Term,
			Data:  log.Data,
			Type:  logType,
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
	output, err := p.adapter.Execute(&domain.ExecuteInput{
		Type: domain.LogType(input.Type),
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}

	var e string
	if output.Err != nil {
		e = output.Err.Error()
	}

	return &p2pRpc.ExecuteOutput{
		Data: output.Out,
		Err:  e,
	}, nil
}

func (p *rpcServerPort) ClusterConfiguration(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.ClusterConfigurationOutput, error) {
	output, err := p.adapter.ClusterConfiguration()
	if err != nil {
		return nil, err
	}

	peers := make(map[string]*p2pRpc.Peer, len(output.Peers))
	for _, peer := range output.Peers {
		peers[peer.Id] = &p2pRpc.Peer{
			Id:     peer.Id,
			Host:   peer.Target(),
			Active: peer.Active,
		}
	}

	return &p2pRpc.ClusterConfigurationOutput{
		Peers:           peers,
		LeaderId:        output.LeaderId,
		ElectionTimeout: uint32(output.ElectionTimeout),
		LeaderHeartbeat: uint32(output.LeaderHeartbeat),
	}, nil
}

func (p *rpcServerPort) Shutdown(ctx context.Context, input *p2pRpc.Nil) (*p2pRpc.Nil, error) {
	p.adapter.Shutdown()
	return &p2pRpc.Nil{}, nil
}
