package rpcsender

import (
	"graft/app/domain/entity"
)

type repository interface {
	AppendEntries(peer entity.Peer, input *entity.AppendEntriesInput) (*entity.AppendEntriesOutput, error)
	RequestVote(peer entity.Peer, input *entity.RequestVoteInput) (*entity.RequestVoteOutput, error)
}

type UseCase interface {
	RunFollower(follower Follower)
	RunCandidate(candadidate Candidate)
	RunLeader(leader Leader)
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

type Follower interface {
	role
	UpgradeCandidate()
}

type Candidate interface {
	role
	downgrader
	broadcaster
	IncrementTerm()
	GetRequestVoteInput() entity.RequestVoteInput
	GetQuorum() int
	UpgradeLeader()
}

type Leader interface {
	role
	downgrader
	broadcaster
	GetAppendEntriesInput() entity.AppendEntriesInput
}
