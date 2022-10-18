package runner

import (
	"graft/src2/entity"
	"log"
	"sync"
	"time"
)

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) RunFollower(follower Follower) {
run:
	for {
		select {
		case <-follower.Timeout().C:
			follower.UpgradeCandidate()
			break run
		default:
			continue
		}
	}
}

func (s *Service) RunCandidate(candidate Candidate) {
	if s.startElection(candidate) {
		return
	}

run:
	for {
		select {
		case <-candidate.Timeout().C:
			if s.startElection(candidate) {
				break run
			}
		default:
			continue
		}
	}
}

func (s *Service) startElection(candidate Candidate) bool {
	log.Println("start election")
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
		return true
	}

	return false
}

func (s *Service) RunLeader(leader Leader) {
	tick := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-tick.C:
			s.sendHeartbeat(leader)
		default:
			continue
		}
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
