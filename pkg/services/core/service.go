package core

import (
	"sync/atomic"

	"graft/pkg/domain"

	log "github.com/sirupsen/logrus"
)

type synchronise struct {
	running    binaryLock
	saving     binaryLock
	committing binaryLock
}

type service struct {
	node    *domain.Node
	repo    repository
	persist persister

	timeout *timeout
	sync    synchronise
}

func NewService(
	node *domain.Node,
	repository repository,
	persister persister,
) *service {
	config := node.GetClusterConfiguration()
	timeout := newTimeout(config.ElectionTimeout, config.LeaderHeartbeat)
	return &service{
		repo:    repository,
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

func (s *service) followerFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		s.node.UpgradeCandidate()
	}
}

func (s *service) candidateFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		if wonElection := s.runElection(); wonElection {
			s.node.UpgradeLeader()
		}
	}
}

func (s *service) leaderFlow() {
	// Prevent multiple flow from running together
	if s.sync.running.Lock() {
		defer s.sync.running.Unlock()

		if quorumReached := s.heartbeat(); !quorumReached {
			log.Debug("STEP DOWN")
			s.node.DowngradeFollower(s.node.CurrentTerm())
		}
	}
}

func (s *service) runFollower() {
	for s.node.Role() == domain.Follower {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.followerFlow()

		case <-s.node.Commit:
			go s.commit()

		case <-s.node.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.node.SaveState:
			go s.saveState()

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) runCandidate() {
	for s.node.Role() == domain.Candidate {
		select {
		case <-s.timeout.ElectionTimer.C:
			go s.candidateFlow()

		case <-s.node.ResetElectionTimer:
			go s.timeout.resetElectionTimer()

		case <-s.node.Commit:
			go s.commit()

		case <-s.node.SaveState:
			go s.saveState()

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) runLeader() {
	for s.node.Role() == domain.Leader {
		select {
		case <-s.timeout.LeaderTicker.C:
			go s.leaderFlow()

		// Should batch execution
		case <-s.node.SynchronizeLogs:
			go (func() {
				s.timeout.resetLeaderTicker()
				s.leaderFlow()
			})()

		case <-s.node.Commit:
			go s.commit()

		case <-s.node.SaveState:
			go s.saveState()

		case <-s.node.ResetLeaderTicker:
			go s.timeout.resetLeaderTicker()

		case <-s.node.ShiftRole:
			return
		}
	}
}

func (s *service) commit() {
	if s.sync.committing.Lock() {
		defer s.sync.committing.Unlock()
		s.node.ApplyLogs()
	}
}

func (s *service) saveState() {
	if s.sync.saving.Lock() {
		defer s.sync.saving.Unlock()
		s.persist.Save(s.node.ToPersistent())
	}
}

func (s *service) runElection() bool {
	// Prevent candidate with outdated logs
	// From running election
	if !s.preVote() {
		log.Debug("FAILED PRE-VOTE")
		return false
	}

	s.node.IncrementCandidateTerm()

	return s.requestVote()
}

func (s *service) preVote() bool {
	input := s.node.RequestVoteInput()
	quorum := s.node.Quorum()
	var prevotesGranted uint32 = 1 // vote for self

	preVoteRoutine := func(p domain.Peer) {
		if res, err := s.repo.PreVote(p, &input); err == nil {
			if res.VoteGranted {
				atomic.AddUint32(&prevotesGranted, 1)
			}
		}
	}
	s.node.Broadcast(preVoteRoutine, domain.BroadcastActive)

	quorumReached := int(prevotesGranted) >= quorum
	return quorumReached
}

func (s *service) requestVote() bool {
	input := s.node.RequestVoteInput()
	quorum := s.node.Quorum()
	var votesGranted uint32 = 1 // vote for self

	gatherVotesRoutine := func(p domain.Peer) {
		if res, err := s.repo.RequestVote(p, &input); err == nil {
			if res.Term > s.node.CurrentTerm() {
				s.node.DowngradeFollower(res.Term)
				return
			}
			if res.VoteGranted {
				atomic.AddUint32(&votesGranted, 1)
			}
		}
	}
	s.node.Broadcast(gatherVotesRoutine, domain.BroadcastActive)

	isCandidate := s.node.Role() == domain.Candidate
	quorumReached := int(votesGranted) >= quorum
	return isCandidate && quorumReached
}

func (s *service) heartbeat() bool {
	quorumReached := s.synchronizeLogs()

	if s.node.Role() == domain.Leader {
		s.node.UpdateNewCommitIndex()
	}

	return quorumReached
}

func (s *service) synchronizeLogs() bool {
	quorum := s.node.Quorum()
	var peersAlive uint32 = 1 // self

	synchroniseLogsRoutine := func(p domain.Peer) {
		input := s.node.AppendEntriesInput(p.Id)
		if res, err := s.repo.AppendEntries(p, &input); err == nil {
			atomic.AddUint32(&peersAlive, 1)
			if res.Term > s.node.CurrentTerm() {
				s.node.DowngradeFollower(res.Term)
				return
			}
			if res.Success {
				// newPeerLastLogIndex is always the sum of len(entries) and prevLogIndex
				// Even if peer was send logs it already has, because prevLogIndex is the number
				// of logs already contained at least in peer, and len(entries) is the additional
				// entries accepted.
				newPeerLastLogIndex := input.PrevLogIndex + uint32(len(input.Entries))
				s.node.SetNextMatchIndex(p.Id, newPeerLastLogIndex)
			} else {
				s.node.DecrementNextIndex(p.Id)
			}
		}
	}
	s.node.Broadcast(synchroniseLogsRoutine, domain.BroadcastAll)

	quorumReached := int(peersAlive) >= quorum
	return quorumReached
}
