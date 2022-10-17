package roles

import (
	"sync"
	"time"
)

type Broadcaster interface {
	Broadcast()
}

// Candidate and Leader are broadcaster
// because they broadcast to the peers

type Voter interface {
	Vote()
	CanVote()
}

// Follower is a voter

// This is wrong, because there wont be multiple struct of candidate, follower or leader
// Candidate, follower and leader should be struct
// type Candidate interface {
// 	StartElection()
// }
// type Follower interface {
// 	RaiseToCandidate()
// }
// type Leader interface {
// 	Broadcast()
// }

type Role struct{}

type Base struct {
	HeartbeatTicker *time.Ticker
	ElectionTimer   *time.Timer
	Mu              sync.Mutex
}
type FollowerState struct {
	CommitIndex uint16
	LastApplied uint16
}

type LeaderState struct {
	FollowerState
	NextIndex  []string
	MatchIndex []string
}

func (r *Follower) RaiseToCandidate() {}
func (r *Follower) Vote()             {}
func (r *Follower) StartElection()    {}

func (r *Candidate) PromoteToLeader()     {}
func (r *Candidate) DowngradeToFollower() {}
func (r *Candidate) StartElection()       {}

func (r *Leader) DowngradeToFollower() {}
func (r *Leader) Broadcast()           {}

// Receive external request
type RpcSever struct{}

// RequestVote / AppendEntries
