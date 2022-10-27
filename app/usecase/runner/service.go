package runner

import (
	"graft/app/domain/entity"
	srv "graft/app/domain/service"
	"sync"
)

type service struct {
	server     *srv.Server
	timeout    *entity.Timeout
	persister  persister
	repository repository
}

func NewService(s *srv.Server, t *entity.Timeout, r repository, p persister) *service {
	return &service{
		repository: r,
		server:     s,
		persister:  p,
		timeout:    t,
	}
}

func (s *service) Run() {
	srv := s.server

	select {
	case <-srv.ShiftRole:
		switch {
		case srv.IsFollower():
			go s.runFollower(srv)
		case srv.IsCandidate():
			go s.runCandidate(srv)
		case srv.IsLeader():
			go s.runLeader(srv)
		}

	case <-srv.SaveState:
		state := srv.GetState()
		s.persister.Save(&state.Persistent)

	case <-srv.ResetElectionTimer:
		s.timeout.ResetElectionTimer()

	case <-srv.ResetLeaderTicker:
		s.timeout.ResetLeaderTicker()
	}
}

func (s *service) runFollower(f follower) {
	for range s.timeout.ElectionTimer.C {
		f.UpgradeCandidate()
		return
	}
}

func (s *service) runCandidate(c candidate) {
	go s.startElection(c)

	for range s.timeout.ElectionTimer.C {
		s.startElection(c)
	}

}

func (s *service) runLeader(l leader) {
	for range s.timeout.LeaderTicker.C {
		s.sendHeartbeat(l)
	}
}

func (s *service) startElection(c candidate) {
	c.IncrementTerm()

	state := c.GetState()
	input := c.GetRequestVoteInput()
	quorum := c.GetQuorum()
	votesGranted := 1 // vote for self

	var m sync.Mutex

	gatherVote := func(p entity.Peer) {
		if res, err := s.repository.RequestVote(p, &input); err == nil {
			if res.Term > state.CurrentTerm {
				c.DowngradeFollower(res.Term, p.Id)
				return
			}
			if res.VoteGranted {
				m.Lock()
				defer m.Unlock()
				votesGranted += 1
			}
		}
	}

	c.Broadcast(gatherVote)

	if votesGranted >= quorum {
		c.UpgradeLeader()
	}
}

func (s *service) sendHeartbeat(l leader) {
	state := l.GetState()
	input := l.GetAppendEntriesInput()

	heartbeat := func(p entity.Peer) {
		if res, err := s.repository.AppendEntries(p, &input); err == nil {
			if res.Term > state.CurrentTerm {
				l.DowngradeFollower(res.Term, p.Id)
				return
			}
		}
	}

	l.Broadcast(heartbeat)
}
