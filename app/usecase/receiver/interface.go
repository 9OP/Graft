package receiver

import (
	"graft/app/domain/entity"
)

type repository interface {
	Heartbeat()
	GetState() *entity.FsmState
	SetClusterLeader(leaderId string)
	SetCommitIndex(ind uint32)
	DowngradeFollower(term uint32)
	GrantVote(id string, lastLogIndex uint32, lastLogTerm uint32) bool
	DeleteLogsFrom(index uint32)
	AppendLogs(entries []entity.LogEntry, prevLogIndex uint32)
}

type UseCase interface {
	AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}
