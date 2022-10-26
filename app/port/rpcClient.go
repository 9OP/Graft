package port

import (
	"fmt"
	"graft/app/adapter"
	"graft/app/adapter/rpc"
	"graft/app/domain/entity"
)

// Repository in use case
//
// type RpcClientPort interface {
// 	AppendEntries(peer entity.Peer, input entity.AppendEntriesInput) (entity.AppendEntriesOutput, error)
// 	RequestVote(peer entity.Peer, input entity.RequestVoteInput) (entity.RequestVoteOutput, error)
// }

type RpcClientPort struct{}

func (p *RpcClientPort) AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error) {
	rpcClient := adapter.RpcClient{}

	target := fmt.Sprintf("%s:%s", peer.Host, peer.Port)
	output, err := rpcClient.AppendEntries(target, &rpc.AppendEntriesInput{
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

/*
Port implementation needs to:
- inject use case
- inject adapter
- transform adpter input to entity
- runs entity trhough use case
- transform use case respose to adapter output
- returns
*/
