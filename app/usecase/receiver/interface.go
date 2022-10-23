package receiver

import (
	"context"
	"graft/app/entity"
	"graft/app/rpc"
)

type Repository interface {
	Heartbeat()
	DowngradeFollower(term uint32, leaderId string)
	SetClusterLeader(leaderId string)
	GetState() *entity.State
	SetCommitIndex(ind uint32)
	GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool
	DeleteLogsFrom(index uint32)
	AppendLogs(entries []string)
}

type UseCase interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}
