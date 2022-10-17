package receiver

import (
	"context"
	"graft/src2/entity"
	"graft/src2/rpc"
)

type Server interface {
	GetState() entity.State
	Heartbeat()
	GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool
	DowngradeFollower(term uint32)
}

type Repository interface {
	Server
}

type UseCase interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}
