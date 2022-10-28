package runner

import (
	"graft/app/domain/entity"
	srv "graft/app/domain/service"
	"log"
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

	var commitMu *sync.Mutex

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

	case <-srv.Commit:
		log.Println("commit index")
		commitMu.Lock()

		// Move this function to domain because not require repository
		applyLogs := func(mu *sync.Mutex) {
			defer mu.Unlock()
			for srv.GetCommitIndex() > srv.GetState().LastApplied {
				srv.IncrementLastApplied()
				lastApplied := srv.GetState().LastApplied
				srv.ApplyLog(lastApplied)
			}
		}
		go applyLogs(commitMu)
		/*
			while srv.commitIndex > srv.lastLogApplied:
				increment lastLogApplied + 1
				and
				apply logs[lastLogApplied]

			WARNING:
			should not apply the same log multiple times (because 2 <-Commit are received simultaneously)
		*/

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

	// TODO synchronise followers with leader logs:

	heartbeat := func(p entity.Peer) {
		// Get the appropriate entries for the Peer, based on its nextIndex
		entries := []string{}
		lastLogIndex := state.GetLastLogIndex()
		followerLogIndex := state.NextIndex[p.Id]
		for lastLogIndex >= followerLogIndex {
			entries = append(entries, state.GetLogByIndex(followerLogIndex).Value)
			followerLogIndex += 1
		}

		input := l.GetAppendEntriesInput(entries)

		if res, err := s.repository.AppendEntries(p, &input); err == nil {
			if res.Term > state.CurrentTerm {
				l.DowngradeFollower(res.Term, p.Id)
				return
			}

			if !res.Success {
				//
				log.Println("failed")
			}

			if res.Success {
				//
				log.Println("success")
			}
		}
	}

	l.Broadcast(heartbeat)
}
