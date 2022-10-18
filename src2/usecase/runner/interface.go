package runner

import (
	"graft/src2/entity"
	"graft/src2/rpc"
)

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
	GetState() entity.State
}

type timeout interface {
	GetTimeout() *entity.Timeout
}

type broadcaster interface {
	Broadcast(fn func(peer entity.Peer))
}

type downgrader interface {
	DowngradeFollower(term uint32)
}

type Follower interface {
	role
	timeout
	UpgradeCandidate()
}

type Candidate interface {
	role
	timeout
	downgrader
	broadcaster
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
