package runner

import (
	"graft/src2/entity"
	"graft/src2/rpc"
	"time"
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

type Follower interface {
	Timeout() <-chan time.Time
	UpgradeCandidate()
}

type Candidate interface {
	DowngradeFollower(term uint32)
	UpgradeLeader()
	Broadcast(fn func(peer entity.Peer))
}

type Leader interface {
	DowngradeFollower(term uint32)
	Broadcast(fn func(peer entity.Peer))
}
