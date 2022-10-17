package raft_server

import "graft/src2/entity"

type Repository interface {
	GetState() entity.State
}

type UseCase interface {
	RunFollower(follower Follower)
	RunCandidate(candidate Candidate)
	RunLeader(leader Leader)
}

type Follower interface {
	// investigate if property can be private
	UpgradeCandidate()
	GrantVote() bool
}

type Candidate interface {
	DowngradeFollower(term uint32)
	UpgradeLeader()
	GrantVote() bool
}

type Leader interface {
	DowngradeFollower(term uint32)
	Broadcast()
}
