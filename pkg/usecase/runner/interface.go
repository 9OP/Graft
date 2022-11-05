package runner

import (
	"graft/pkg/domain/entity"
)

type repository interface {
	AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}

type persister interface {
	Load() (*entity.PersistentState, error)
	Save(currentTerm uint32, votedFor string, machineLogs []entity.LogEntry) error
}

type UseCase interface {
	Run()
}

type role interface {
	GetState() entity.NodeState
}
type broadcaster interface {
	Broadcast(fn func(peer entity.Peer))
	Quorum() int
}
type downgrader interface {
	DowngradeFollower(term uint32)
}

type follower interface {
	role
	UpgradeCandidate()
}

type candidate interface {
	role
	downgrader
	broadcaster
	IncrementCandidateTerm()
	RequestVoteInput() entity.RequestVoteInput
	UpgradeLeader()
}

type leader interface {
	role
	downgrader
	broadcaster
	DecrementNextIndex(peerId string)
	SetNextMatchIndex(peerId string, index uint32)
	ComputeNewCommitIndex() uint32
	SetCommitIndex(commitIndex uint32)
	AppendEntriesInput(peerId string) entity.AppendEntriesInput
}
