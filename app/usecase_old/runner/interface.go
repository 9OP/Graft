package runner

import (
	"graft/app/domain/entity"
	"graft/app/rpc"
)

// Rename Rpc
type Repository interface {
	AppendEntries(peer entity.Peer, input *rpc.AppendEntriesInput) (*rpc.AppendEntriesOutput, error)
	RequestVote(peer entity.Peer, input *rpc.RequestVoteInput) (*rpc.RequestVoteOutput, error)
}

type UseCase interface {
	RunFollower(follower Follower)
	RunCandidate(candadidate Candidate)
	RunLeader(leader Leader)
}

type role interface {
	GetState() *entity.State
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
	RequestVoteInput() *rpc.RequestVoteInput
	GetQuorum() int
	UpgradeLeader()
}

type Leader interface {
	role
	downgrader
	broadcaster
	AppendEntriesInput() *rpc.AppendEntriesInput
}
