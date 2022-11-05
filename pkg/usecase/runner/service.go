package runner

import (
	"graft/pkg/domain/entity"
	domain "graft/pkg/domain/service"
	"sync"

	log "github.com/sirupsen/logrus"
)

type synchronise struct {
	save      sync.Mutex
	commit    sync.Mutex
	election  sync.Mutex
	heartbeat sync.Mutex
}

type service struct {
	clusterNode *domain.ClusterNode
	timeout     *entity.Timeout
	persister   persister
	repository  repository
	sync        synchronise
}

func NewService(clusterNode *domain.ClusterNode, timeout *entity.Timeout, repository repository, persister persister) *service {
	return &service{
		repository:  repository,
		timeout:     timeout,
		clusterNode: clusterNode,
		persister:   persister,
	}
}

func (s *service) Run() {
	select {
	case role := <-s.clusterNode.ShiftRole:
		go s.runAs(role)

	case <-s.clusterNode.SaveState:
		go s.saveState()

	case <-s.clusterNode.ResetElectionTimer:
		go s.timeout.ResetElectionTimer()

	case <-s.clusterNode.ResetLeaderTicker:
		go s.timeout.ResetLeaderTicker()

	case <-s.clusterNode.Commit:
		go s.commit()
	}

}

func (s *service) runAs(role entity.Role) {
	switch role {
	case entity.Follower:
		s.runFollower(s.clusterNode)
	case entity.Candidate:
		s.runCandidate(s.clusterNode)
	case entity.Leader:
		s.runLeader(s.clusterNode)
	}

}

func (s *service) commit() {
	s.sync.commit.Lock()
	defer s.sync.commit.Unlock()
	s.clusterNode.ApplyLogs()
}

func (s *service) saveState() {
	s.sync.save.Lock()
	defer s.sync.save.Unlock()
	s.persister.Save(
		s.clusterNode.CurrentTerm(),
		s.clusterNode.VotedFor(),
		s.clusterNode.MachineLogs(),
	)
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
				return
			}
		}
	}
}

func (s *service) runLeader(l leader) {
	for range s.timeout.LeaderTicker.C {
		if !s.sendHeartbeat(l) {
			log.Debug("LEADER STEP DOWN")
			term := l.GetState().CurrentTerm()
			l.DowngradeFollower(term)
			return
		}
	}
}

func (s *service) startElection(c candidate) bool {
	s.sync.election.Lock()
	defer s.sync.election.Unlock()

	c.IncrementCandidateTerm()
	state := c.GetState()
	input := c.RequestVoteInput()
	quorum := c.Quorum()
	votesGranted := 1 // vote for self
	var m sync.Mutex

	gatherVotesRoutine := func(p entity.Peer) {
		if res, err := s.repository.RequestVote(p, &input); err == nil {
			if res.Term > state.CurrentTerm() {
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

func (s *service) sendHeartbeat(l leader) bool {
	s.sync.heartbeat.Lock()
	defer s.sync.heartbeat.Unlock()

	state := l.GetState()
	quorum := l.Quorum()
	peersAlive := 1 // self
	var m sync.Mutex

	synchroniseLogsRoutine := func(p entity.Peer) {
		input := l.AppendEntriesInput(p.Id)

		if res, err := s.repository.AppendEntries(p, &input); err == nil {
			if res.Term > state.CurrentTerm() {
				l.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				leaderLastLogIndex := state.LastLogIndex()
				peerNextIndex := state.NextIndexForPeer(p.Id)
				peerMatchIndex := state.MatchIndexForPeer(p.Id)
				shouldUpdate :=
					peerNextIndex != leaderLastLogIndex ||
						peerMatchIndex != leaderLastLogIndex
				if shouldUpdate {
					l.SetNextMatchIndex(p.Id, leaderLastLogIndex)
				}
			} else {
				l.DecrementNextIndex(p.Id)
			}
			m.Lock()
			defer m.Unlock()
			peersAlive += 1
		}
	}
	l.Broadcast(synchroniseLogsRoutine)

	newCommitIndex := l.ComputeNewCommitIndex()
	l.SetCommitIndex(newCommitIndex)
	return peersAlive >= quorum
}
