package runner

import "time"

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
}

type Leader interface {
	DowngradeFollower(term uint32)
	Broadcast()
}
