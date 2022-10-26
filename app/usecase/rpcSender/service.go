package rpcsender

import (
	"fmt"
	"graft/app/domain/entity"
	"sync"
)

type Service struct {
	repo    repository
	timeout *entity.Timeout
}

func NewService(repo repository, timeout *entity.Timeout) *Service {
	return &Service{repo, timeout}
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
	input := candidate.GetRequestVoteInput()
	quorum := candidate.GetQuorum()
	votesGranted := 1 // vote for self

	var m sync.Mutex

	gatherVote := func(p entity.Peer) {
		if res, err := s.repo.RequestVote(p, &input); err == nil {
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

	fmt.Println(votesGranted, quorum)

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
	input := leader.GetAppendEntriesInput()

	heartbeat := func(p entity.Peer) {
		if res, err := s.repo.AppendEntries(p, &input); err == nil {
			if res.Term > state.CurrentTerm {
				leader.DowngradeFollower(res.Term, p.Id)
				return
			}
		}
	}

	leader.Broadcast(heartbeat)
}
