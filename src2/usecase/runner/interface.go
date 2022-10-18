package runner

import (
	"graft/src2/entity"
	"time"
)

type Repository interface {
	AppendEntries()
	RequestVote()
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
	Broadcast(fn func(peer entity.Peer, cd Candidate))
}

type Leader interface {
	DowngradeFollower(term uint32)
	Broadcast(fn func(peer entity.Peer, ld Leader))
}
