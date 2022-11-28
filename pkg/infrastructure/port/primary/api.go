package primaryPort

import (
	"context"

	"graft/pkg/domain"
	"graft/pkg/infrastructure/adapter/clusterRpc"
	"graft/pkg/services/rpc"
)

type ServerAdapter interface {
	// P2P
	AppendEntries(ctx context.Context, input *clusterRpc.AppendEntriesInput) (*clusterRpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error)
	PreVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error)
	// Cluster
	Configuration(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.ConfigurationOutput, error)
	LeadershipTransfer(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
	Execute(ctx context.Context, input *clusterRpc.ExecuteInput) (*clusterRpc.ExecuteOutput, error)
	Shutdown(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
	Ping(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
}

type rpcServerPort struct {
	adapter rpc.UseCase
}

func NewRpcServerPort(adapter rpc.UseCase) *rpcServerPort {
	return &rpcServerPort{adapter}
}

func (p *rpcServerPort) AppendEntries(ctx context.Context, input *clusterRpc.AppendEntriesInput) (*clusterRpc.AppendEntriesOutput, error) {
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

	return &clusterRpc.AppendEntriesOutput{
		Term:    output.Term,
		Success: output.Success,
	}, nil
}

func (p *rpcServerPort) RequestVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	output, err := p.adapter.RequestVote(&domain.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		return nil, err
	}

	return &clusterRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}

func (p *rpcServerPort) PreVote(ctx context.Context, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error) {
	output, err := p.adapter.PreVote(&domain.RequestVoteInput{
		CandidateId:  input.CandidateId,
		Term:         input.Term,
		LastLogIndex: input.LastLogIndex,
		LastLogTerm:  input.LastLogTerm,
	})
	if err != nil {
		return nil, err
	}

	return &clusterRpc.RequestVoteOutput{
		Term:        output.Term,
		VoteGranted: output.VoteGranted,
	}, nil
}

func (p *rpcServerPort) Execute(ctx context.Context, input *clusterRpc.ExecuteInput) (*clusterRpc.ExecuteOutput, error) {
	output, err := p.adapter.Execute(&domain.ExecuteInput{
		Type: domain.LogType(input.Type),
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}

	var outputErr string
	if output.Err != nil {
		outputErr = output.Err.Error()
	}

	return &clusterRpc.ExecuteOutput{
		Data: output.Out,
		Err:  outputErr,
	}, nil
}

func (p *rpcServerPort) LeadershipTransfer(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	p.adapter.LeadershipTransfer()
	return &clusterRpc.Nil{}, nil
}

func (p *rpcServerPort) Configuration(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.ConfigurationOutput, error) {
	output, err := p.adapter.Configuration()
	if err != nil {
		return nil, err
	}

	peers := make(map[string]*clusterRpc.Peer, len(output.Peers))
	for _, peer := range output.Peers {
		peers[peer.Id] = &clusterRpc.Peer{
			Id:     peer.Id,
			Host:   peer.Target(),
			Active: peer.Active,
		}
	}

	return &clusterRpc.ConfigurationOutput{
		Peers:           peers,
		LeaderId:        output.LeaderId,
		Fsm:             output.Fsm,
		ElectionTimeout: uint32(output.ElectionTimeout),
		LeaderHeartbeat: uint32(output.LeaderHeartbeat),
	}, nil
}

func (p *rpcServerPort) Shutdown(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	p.adapter.Shutdown()
	return &clusterRpc.Nil{}, nil
}

func (p *rpcServerPort) Ping(ctx context.Context, input *clusterRpc.Nil) (*clusterRpc.Nil, error) {
	err := p.adapter.Ping()
	return &clusterRpc.Nil{}, err
}
