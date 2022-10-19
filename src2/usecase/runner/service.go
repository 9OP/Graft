package runner

import (
	"graft/src2/entity"
	"graft/src2/rpc"
	"sync"
)

type Service struct {
	repository Repository
	timeout    *entity.Timeout
	ticker     *entity.Ticker
}

func NewService(repo Repository, timeout *entity.Timeout, ticker *entity.Ticker) *Service {
	return &Service{repository: repo, timeout: timeout, ticker: ticker}
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
			s.startElection(candidate, signal)
		case <-signal:
			break run
		}
	}

}

func (s *Service) startElection(candidate Candidate, signal chan struct{}) {
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
	s.ticker.Start()

	for range s.ticker.C {
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
