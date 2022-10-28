package runner

import (
	"graft/app/domain/entity"
)

type repository interface {
	AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}

type persister interface {
	Load() (*entity.Persistent, error)
	Save(state *entity.Persistent) error
}

type UseCase interface {
	Run()
}

type role interface {
	GetState() entity.FsmState
}
type broadcaster interface {
	Broadcast(fn func(peer entity.Peer))
}
type downgrader interface {
	DowngradeFollower(term uint32, leaderId string)
}

type follower interface {
	role
	UpgradeCandidate()
}

type candidate interface {
	role
	downgrader
	broadcaster
	IncrementTerm()
	GetRequestVoteInput() entity.RequestVoteInput
	GetQuorum() int
	UpgradeLeader()
}

type leader interface {
	role
	downgrader
	broadcaster
	GetAppendEntriesInput(entries []string) entity.AppendEntriesInput
}
