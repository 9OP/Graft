package runner

import (
	"graft/src2/entity"
	"sync"
	"time"
)

type Service struct {
	repository Repository
	timeout    *entity.Timeout
}

func NewService(repo Repository, timeout *entity.Timeout) *Service {
	return &Service{repository: repo, timeout: timeout}
}

func (s *Service) RunFollower(follower Follower) {
	for range s.timeout.C {
		follower.UpgradeCandidate()
		return
	}
}

func (s *Service) RunCandidate(candidate Candidate) {
	signal := make(chan struct{}, 1)

	go s.startElection(candidate, signal)

run:
	for {
		select {
		case <-s.timeout.C:
			go s.restartElection(candidate, signal)
		case <-signal:
			break run
		}
	}

}

func (s *Service) restartElection(candidate Candidate, signal chan struct{}) {
	candidate.IncrementTerm()
	s.startElection(candidate, signal)
}

func (s *Service) startElection(candidate Candidate, signal chan struct{}) {
	state := candidate.GetState()
	input := candidate.RequestVoteInput()
	quorum := candidate.GetQuorum()
	votesGranted := 1 // vote for self
	var m sync.Mutex

	gatherVote := func(p entity.Peer) {
		if res, err := s.repository.RequestVote(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				candidate.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				m.Lock()
				defer m.Unlock()
				votesGranted += 1
			}
		}
	}

	candidate.Broadcast(gatherVote)

	if votesGranted >= quorum {
		candidate.UpgradeLeader()
		signal <- struct{}{}
	}
}

func (s *Service) RunLeader(leader Leader) {
	tick := time.NewTicker(50 * time.Millisecond)

	for range tick.C {
		s.sendHeartbeat(leader)
	}
}

func (s *Service) sendHeartbeat(leader Leader) {
	state := leader.GetState()
	input := leader.AppendEntriesInput()

	heartbeat := func(p entity.Peer) {
		if res, err := s.repository.AppendEntries(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				leader.DowngradeFollower(res.Term)
				return
			}
		}
	}

	leader.Broadcast(heartbeat)
}
