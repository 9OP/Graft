package secondaryPort

import (
	"errors"
	"net/netip"

	"graft/pkg/domain"
	"graft/pkg/infrastructure/adapter/clusterRpc"
)

type ClientAdapter interface {
	AppendEntries(target string, input *clusterRpc.AppendEntriesInput) (*clusterRpc.AppendEntriesOutput, error)
	RequestVote(target string, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error)
	PreVote(target string, input *clusterRpc.RequestVoteInput) (*clusterRpc.RequestVoteOutput, error)
	//
	Execute(target string, input *clusterRpc.ExecuteInput) (*clusterRpc.ExecuteOutput, error)
	LeadershipTransfer(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
	Configuration(target string, input *clusterRpc.Nil) (*clusterRpc.ConfigurationOutput, error)
	Shutdown(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
	Ping(target string, input *clusterRpc.Nil) (*clusterRpc.Nil, error)
}

type rpcClientPort struct {
	adapter ClientAdapter
}

func NewRpcClientPort(adapter ClientAdapter) *rpcClientPort {
	return &rpcClientPort{adapter}
}

func (p *rpcClientPort) AppendEntries(peer domain.Peer, input *domain.AppendEntriesInput) (*domain.AppendEntriesOutput, error) {
	entries := make([]*clusterRpc.LogEntry, 0, len(input.Entries))
	for _, log := range input.Entries {
		logType := clusterRpc.LogType(clusterRpc.LogType_value[log.Type.String()])
		entry := &clusterRpc.LogEntry{
			Index: log.Index,
			Term:  log.Term,
			Data:  log.Data,
			Type:  logType,
		}
		entries = append(entries, entry)
	}

	output, err := p.adapter.AppendEntries(peer.Target(), &clusterRpc.AppendEntriesInput{
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
	output, err := p.adapter.RequestVote(peer.Target(), &clusterRpc.RequestVoteInput{
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
	output, err := p.adapter.PreVote(peer.Target(), &clusterRpc.RequestVoteInput{
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

func (p *rpcClientPort) Execute(peer domain.Peer, input *domain.ExecuteInput) (*domain.ExecuteOutput, error) {
	output, err := p.adapter.Execute(peer.Target(), &clusterRpc.ExecuteInput{
		Type: clusterRpc.LogType(input.Type),
		Data: input.Data,
	})
	if err != nil {
		return nil, err
	}

	return &domain.ExecuteOutput{
		Out: output.Data,
		Err: errors.New(output.Err),
	}, nil
}

func (p *rpcClientPort) LeadershipTransfer(peer domain.Peer) error {
	_, err := p.adapter.LeadershipTransfer(peer.Target(), &clusterRpc.Nil{})
	return err
}

func (p *rpcClientPort) Configuration(peer domain.Peer) (*domain.ClusterConfiguration, error) {
	output, err := p.adapter.Configuration(peer.Target(), &clusterRpc.Nil{})
	if err != nil {
		return nil, err
	}

	peers := make(domain.Peers, len(output.Peers))
	for _, peer := range output.Peers {
		host, _ := netip.ParseAddrPort(peer.Host)
		peers[peer.Id] = domain.Peer{
			Id:     peer.Id,
			Host:   host,
			Active: peer.Active,
		}
	}

	return &domain.ClusterConfiguration{
		Peers:           peers,
		LeaderId:        output.LeaderId,
		ElectionTimeout: int(output.ElectionTimeout),
		LeaderHeartbeat: int(output.LeaderHeartbeat),
	}, nil
}

func (p *rpcClientPort) Shutdown(peer domain.Peer) error {
	_, err := p.adapter.Shutdown(peer.Target(), &clusterRpc.Nil{})
	return err
}

func (p *rpcClientPort) Ping(peer domain.Peer) error {
	_, err := p.adapter.Ping(peer.Target(), &clusterRpc.Nil{})
	return err
}
