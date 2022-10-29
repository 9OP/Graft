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

func NewService(s *srvc.Server, t *entity.Timeout, r repository, p persister) *service {
	return &service{
		repository: r,
		server:     s,
		persister:  p,
		timeout:    t,
	}
}

func (s *service) Run() {
	events := s.server.Signals

	select {
	case role := <-events.ShiftRole:
		go s.runAs(role)

	case <-events.SaveState:
		go s.saveState()

	case <-events.ResetElectionTimer:
		go s.timeout.ResetElectionTimer()

	case <-events.ResetLeaderTicker:
		go s.timeout.ResetLeaderTicker()

	case <-events.Commit:
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
	s.persister.Save(&state.Persistent)
}

func (s *service) runFollower(f follower) {
	for range s.timeout.ElectionTimer.C {
		f.UpgradeCandidate()
		return
	}
}

func (s *service) runCandidate(c candidate) {
	// Start election when becomming candidate
	if s.startElection(c) {
		c.UpgradeLeader()
		return
	}

	// When timeout, restart election
	for range s.timeout.ElectionTimer.C {
		if s.startElection(c) {
			c.UpgradeLeader()
			return
		}
	}
}

func (s *service) runLeader(l leader) {
	for range s.timeout.LeaderTicker.C {
		s.sendHeartbeat(l)
	}
}

func (s *service) startElection(c candidate) bool {
	// Prevent starting multiple election
	s.synchronise.election.Lock()
	defer s.synchronise.election.Unlock()

	c.IncrementTerm()

	state := c.GetState()
	input := c.GetRequestVoteInput()
	quorum := c.GetQuorum()
	votesGranted := 1 // vote for self

	var m sync.Mutex
	gatherVotesRoutine := func(p entity.Peer) {
		if res, err := s.repository.RequestVote(p, &input); err == nil {
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
	// Prevent starting multiple heartbeat
	s.synchronise.heartbeat.Lock()
	defer s.synchronise.heartbeat.Unlock()

	state := l.GetState()

	synchroniseLogsRoutine := func(p entity.Peer) {
		// entries := l.GetPeerNewEntries(p.Id)
		entries := []string{}
		input := l.GetAppendEntriesInput(entries)
		if res, err := s.repository.AppendEntries(p, &input); err == nil {
			if res.Term > state.CurrentTerm {
				l.DowngradeFollower(res.Term)
				return
			}
			// if res.Success {
			// 	leaderLastLogIndex := state.GetLastLogIndex()
			// 	l.SetNextIndex(p.Id, leaderLastLogIndex)
			// 	l.SetMatchIndex(p.Id, leaderLastLogIndex)
			// } else {
			// 	l.DecrementNextIndex(p.Id)
			// }
		}
	}
	l.Broadcast(synchroniseLogsRoutine)

	/*
		If there exists an N such that:
			- N > commitIndex,
			- a majority of matchIndex[i] ≥ N
			- log[N].term == currentTerm:
		then set commitIndex = N
	*/

	// // Refresh state
	// state = l.GetState()
	// // Upper value of N
	// lastLogIndex := state.GetLastLogIndex()
	// // Lower value of N
	// commitIndex := state.CommitIndex
	// // Look for N, starting from latest log
	// quorum := s.server.GetQuorum()
	// for N := lastLogIndex; N > commitIndex; N-- {
	// 	// Get a majority for which matchIndex >= n
	// 	count := 0
	// 	for _, matchIndex := range state.MatchIndex {
	// 		if matchIndex >= N {
	// 			count += 1
	// 		}
	// 	}
	// 	if count >= quorum {
	// 		log := state.GetLogByIndex(N)
	// 		if log.Term == state.CurrentTerm {
	// 			s.server.SetCommitIndex(N)
	// 			break
	// 		}
	// 	}
	// }

}
