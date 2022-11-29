package core

import (
	"graft/pkg/domain"
	"graft/pkg/services/lib"

	"graft/pkg/utils/log"
)

type service struct {
	node    *domain.Node
	client  lib.Client
	persist persister
	timeout *timeout
	running binaryLock
}

func NewService(
	node *domain.Node,
	client lib.Client,
	persister persister,
) *service {
	config := node.GetClusterConfiguration()
	timeout := newTimeout(config.ElectionTimeout, config.LeaderHeartbeat)
	return &service{
		client:  client,
		timeout: timeout,
		node:    node,
		persist: persister,
	}
}

func (s *service) Run() {
	for {
		switch s.node.Role() {
		case domain.Follower:
			s.runFollower()
		case domain.Candidate:
			s.runCandidate()
		case domain.Leader:
			s.runLeader()
		}
	}
}

func (s *service) followerTimeout() {
	// Prevent multiple flow from running together
	if s.running.Lock() {
		defer s.running.Unlock()

		s.node.UpgradeCandidate()
	}
}

func (s *service) candidateTimeout() {
	// Prevent multiple flow from running together
	if s.running.Lock() {
		defer s.running.Unlock()

		if wonElection := s.runElection(); wonElection {
			s.node.UpgradeLeader()
		}
	}
}

func (s *service) leaderHeartbeat() {
	// Prevent multiple flow from running together
	if s.running.Lock() {
		defer s.running.Unlock()

		if quorumReached := s.heartbeat(); !quorumReached {
			log.Debugf("leader step down")
			s.node.DowngradeFollower(s.node.CurrentTerm())
		}
	}
}

func (s *service) runFollower() {
	for s.node.Role() == domain.Follower {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.followerTimeout()

		case <-s.node.Commit:
			go s.node.ApplyLogs()

		case <-s.node.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.node.SaveState:
			go s.persist.Save(s.node.ToPersistent())

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) runCandidate() {
	for s.node.Role() == domain.Candidate {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.candidateTimeout()

		case <-s.node.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.node.Commit:
			go s.node.ApplyLogs()

		case <-s.node.SaveState:
			go s.persist.Save(s.node.ToPersistent())

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) runLeader() {
	for s.node.Role() == domain.Leader {
		select {
		case <-s.timeout.LeaderTicker.C:
			go s.leaderHeartbeat()

		// Should batch execution
		case <-s.node.SynchronizeLogs:
			go (func() {
				s.timeout.resetLeaderTicker()
				s.leaderHeartbeat()
			})()

		case <-s.node.Commit:
			go s.node.ApplyLogs()

		case <-s.node.SaveState:
			go s.persist.Save(s.node.ToPersistent())

		case <-s.node.ResetLeaderTicker:
			go s.timeout.resetLeaderTicker()

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) runElection() bool {
	if !lib.PreVote(s.node, s.client) {
		s.node.DowngradeFollower(s.node.CurrentTerm())
		return false
	}
	s.node.IncrementCandidateTerm()
	return lib.RequestVote(s.node, s.client)
}

func (s *service) heartbeat() bool {
	quorumReached := lib.AppendEntries(s.node, s.client)
	s.node.ComputeNewCommitIndex()
	return quorumReached
}
