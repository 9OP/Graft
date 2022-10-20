package receiver

import (
	"context"
	"graft/app/entity"
	"graft/app/rpc"
)

type Repository interface {
	GetState() entity.ImmerState
	Heartbeat()
	GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool
	DowngradeFollower(term uint32, leaderId string)
	SetClusterLeader(leaderId string)
	DeleteLogsFrom(n int)
	AppendLogs(entries []string)
	SetCommitIndex(ind uint32)
}

type UseCase interface {
	AppendEntries(ctx context.Context, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(ctx context.Context, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}
