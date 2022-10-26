package port

import (
	"fmt"
	"graft/app/domain/entity"
	"graft/app/infrastructure/adapter/rpc"
)

// Repository in use case runner
//
// type Repo interface {
// 	AppendEntries(peer entity.Peer, input entity.AppendEntriesInput) (entity.AppendEntriesOutput, error)
// 	RequestVote(peer entity.Peer, input entity.RequestVoteInput) (entity.RequestVoteOutput, error)
// }

type RpcClientAdapter interface {
	AppendEntries(target string, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(target string, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}

type rpcClientPort struct {
	adapter RpcClientAdapter
}

func NewRpcClientPort(adapter RpcClientAdapter) *rpcClientPort {
	return &rpcClientPort{adapter}
}

func (p *rpcClientPort) AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	target := fmt.Sprintf("%s:%s", peer.Host, peer.Port)
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
	target := fmt.Sprintf("%s:%s", peer.Host, peer.Port)
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
