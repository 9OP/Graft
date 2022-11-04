package receiver

import (
	"graft/pkg/domain/entity"
)

type repository interface {
	Heartbeat()
	GetState() entity.NodeState
	SetClusterLeader(leaderId string)
	SetCommitIndex(ind uint32)
	DowngradeFollower(term uint32)
	GrantVote(peerId string)
	DeleteLogsFrom(index uint32)
	AppendLogs(prevLogIndex uint32, entries ...entity.LogEntry)
}

type UseCase interface {
	AppendEntries(input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}
