package raft_server

import "graft/src/entity"

type Repository interface {
	GetState() entity.State
}

type UseCase interface {
	RunFollower(follower Follower)
	RunCandidate(candidate Candidate)
	RunLeader(leader Leader)
}

type Follower interface {
}

type Candidate interface {
}

type Leader interface {
}

// // UseCase

// type IFollower interface {
// 	upCandidate()
// 	grantVote()
// }

// type ICandidate interface {
// 	downFollower()
// 	upLeader()
// 	grantVote()
// }

// type ILeader interface {
// 	downFollower()
// 	// broadcast ?
// }
