package runner

import (
	"graft/app/entity"
	entitynew "graft/app/entity_new"
	"graft/app/rpc"
	"sync"
)

type Service struct {
	repository Repository
	timeout    *entitynew.Timeout
	// // Create one object that hold a timer and a ticker
	// timeout *entity.Timeout
	// ticker  *entity.Ticker
}

func NewService(repo Repository, timeout *entitynew.Timeout) *Service {
	return &Service{repository: repo, timeout: timeout}
}

func (s *Service) RunFollower(follower Follower) {
	for range s.timeout.ElectionTimer.C {
		follower.UpgradeCandidate()
		return
	}
}

func (s *Service) RunCandidate(candidate Candidate) {
	go s.startElection(candidate)

	for range s.timeout.ElectionTimer.C {
		s.startElection(candidate)
	}

}

func (s *Service) startElection(candidate Candidate) {
	candidate.IncrementTerm()

	state := candidate.GetState()
	input := candidate.RequestVoteInput()
	quorum := candidate.GetQuorum()
	votesGranted := 1 // vote for self

	var m sync.Mutex

	gatherVote := func(p entity.Peer) {
		var err error
		var res *rpc.RequestVoteOutput
		if res, err = s.repository.RequestVote(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				candidate.DowngradeFollower(res.Term, p.Id)
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
	}
}

func (s *Service) RunLeader(leader Leader) {
	for range s.timeout.LeaderTicker.C {
		s.sendHeartbeat(leader)
	}
}

func (s *Service) sendHeartbeat(leader Leader) {
	state := leader.GetState()
	input := leader.AppendEntriesInput()

	heartbeat := func(p entity.Peer) {
		if res, err := s.repository.AppendEntries(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				leader.DowngradeFollower(res.Term, p.Id)
				return
			}
		}
	}

	leader.Broadcast(heartbeat)
}
