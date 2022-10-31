package runner

import (
	"graft/app/domain/entity"
	srvc "graft/app/domain/service"
	"sync"
)

type synchronise struct {
	commit    sync.Mutex
	election  sync.Mutex
	heartbeat sync.Mutex
}

type service struct {
	server     *srvc.Server
	timeout    *entity.Timeout
	persister  persister
	repository repository
	synchronise
}

// TODO: Move timeout to this file instead of being in the domain
// Timeout is applicative logic, not domain logic
func NewService(s *srvc.Server, t *entity.Timeout, r repository, p persister) *service {
	return &service{
		repository: r,
		server:     s,
		persister:  p,
		timeout:    t,
	}
}

func (s *service) Run() {
	select {
	case role := <-s.server.ShiftRole:
		go s.runAs(role)

	case <-s.server.SaveState:
		go s.saveState()

	case <-s.server.ResetElectionTimer:
		go s.timeout.ResetElectionTimer()

	case <-s.server.ResetLeaderTicker:
		go s.timeout.ResetLeaderTicker()

	case <-s.server.Commit:
		go s.commit()
	}

}

func (s *service) runAs(role entity.Role) {
	switch role {
	case entity.Follower:
		s.runFollower(s.server)
	case entity.Candidate:
		s.runCandidate(s.server)
	case entity.Leader:
		s.runLeader(s.server)
	}

}

func (s *service) commit() {
	s.synchronise.commit.Lock()
	defer s.synchronise.commit.Unlock()
	s.server.ApplyLogs()
}

func (s *service) saveState() {
	state := s.server.GetState()
	s.persister.Save(state.Persistent)
}

func (s *service) runFollower(f follower) {
	for range s.timeout.ElectionTimer.C {
		f.UpgradeCandidate()
		return
	}
}

func (s *service) runCandidate(c candidate) {
	if wonElection := s.startElection(c); wonElection {
		c.UpgradeLeader()
	} else {
		// Restart election until: upgrade / downgrade
		for range s.timeout.ElectionTimer.C {
			if wonElection := s.startElection(c); wonElection {
				c.UpgradeLeader()
				break
			}
		}
	}
}

func (s *service) runLeader(l leader) {
	for range s.timeout.LeaderTicker.C {
		s.sendHeartbeat(l)
	}
}

func (s *service) startElection(c candidate) bool {
	s.synchronise.election.Lock()
	defer s.synchronise.election.Unlock()

	c.IncrementTerm()
	state := c.GetState()
	input := c.GetRequestVoteInput()
	quorum := c.GetQuorum()
	votesGranted := 1 // vote for self

	var m sync.Mutex
	gatherVotesRoutine := func(p entity.Peer) {
		if res, err := s.repository.RequestVote(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				c.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				m.Lock()
				defer m.Unlock()
				votesGranted += 1
			}
		}
	}
	c.Broadcast(gatherVotesRoutine)

	return votesGranted >= quorum
}

func (s *service) sendHeartbeat(l leader) {
	s.synchronise.heartbeat.Lock()
	defer s.synchronise.heartbeat.Unlock()

	state := l.GetState()

	synchroniseLogsRoutine := func(p entity.Peer) {
		input := l.GetAppendEntriesInput(p.Id)

		if res, err := s.repository.AppendEntries(p, input); err == nil {
			if res.Term > state.CurrentTerm {
				l.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				// TODO: optim, set index/match only if different from leaderLastLogIndex
				leaderLastLogIndex := state.GetLastLogIndex()
				l.SetNextIndex(p.Id, leaderLastLogIndex)
				l.SetMatchIndex(p.Id, leaderLastLogIndex)
			} else {
				l.DecrementNextIndex(p.Id)
			}
		}
	}
	l.Broadcast(synchroniseLogsRoutine)

	newCommitIndex := l.ComputeNewCommitIndex()
	l.SetCommitIndex(newCommitIndex)
}
